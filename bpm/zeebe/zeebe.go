package zeebe

import (
	"context"
	"fmt"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/entities"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/worker"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
	"github.com/exluap/kit/bpm"
	"github.com/exluap/kit/log"
	"path/filepath"
	"strings"
	"time"
)

const ERR_EXHAUSTED_RESOURCES_MAX_RETRY = 10

type engineImpl struct {
	params     *bpm.Params
	client     zbc.Client
	jobWorkers []worker.JobWorker
	logger     log.CLoggerFunc
}

func NewEngine(logger log.CLoggerFunc) bpm.Engine {

	zeebe := &engineImpl{
		params:     &bpm.Params{},
		jobWorkers: []worker.JobWorker{},
		logger:     logger,
	}

	return zeebe

}

func (z *engineImpl) l() log.CLogger {
	return z.logger().Cmp("zeebe")
}

func (z *engineImpl) Open(params *bpm.Params) error {

	z.params = params

	// if already opened, close it
	if err := z.Close(); err != nil {
		return ErrZeebeClose(err)
	}

	if z.params.Port == "" || z.params.Host == "" {
		return ErrZeebeConnInvalidParams(z.params.Host, z.params.Port)
	}

	if zc, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         fmt.Sprintf("%s:%s", z.params.Host, z.params.Port),
		UsePlaintextConnection: true,
	}); err == nil {
		z.client = zc
		z.l().Mth("open").Inf("ok")
	} else {
		return ErrZeebeNewClient(err)
	}

	return nil
}

func (z *engineImpl) IsOpened() bool {
	return z.client != nil
}

func (z *engineImpl) Close() error {
	if z.client != nil {
		err := z.client.Close()
		z.client = nil
		if err != nil {
			return ErrZeebeClose(err)
		}
		z.l().Mth("close").Inf("closed")
	}
	return nil
}

func (z *engineImpl) DeployBPMNsAsync(paths []string) {

	for _, p := range paths {

		absPath, _ := filepath.Abs(p)

		go func(path string) {

			l := z.l().Mth("deploy").F(log.FF{"path": path})
			errRetryCount := 0

			for {
				rs, err := z.client.NewDeployProcessCommand().AddResourceFile(path).Send(context.Background())
				if err != nil {
					if strings.Contains(err.Error(), "ResourceExhausted") && errRetryCount <= ERR_EXHAUSTED_RESOURCES_MAX_RETRY {
						time.Sleep(time.Millisecond * 100)
						errRetryCount++
						l.Warn("retry")
					} else {
						l.E(err).Err()
						return
					}
				} else {
					l.Dbg("deployed").TrcF("details: %v", rs)
					return
				}
			}
		}(absPath)

	}

}

func (z *engineImpl) RegisterTaskHandlers(handlers map[string]interface{}) error {

	for task, handlerFunc := range handlers {
		if f, ok := handlerFunc.(func(client worker.JobClient, job entities.Job)); ok {
			z.jobWorkers = append(z.jobWorkers, z.client.NewJobWorker().JobType(task).Handler(f).Open())
		} else {
			return ErrZeebeConnInvalidHandler()
		}
	}

	return nil

}

func (z *engineImpl) StartProcess(processId string, vars map[string]interface{}) (string, error) {

	l := z.l().Mth("start").F(log.FF{"process-id": processId})

	p, err := z.client.NewCreateInstanceCommand().BPMNProcessId(processId).LatestVersion().VariablesFromMap(vars)
	if err != nil {
		return "", err
	}

	ctx := context.Background()

	pRs, err := p.Send(ctx)
	if err != nil {
		return "", ErrZeebeStartProcess(err, processId)
	}

	l.Dbg("started").TrcF("vars=%v details=%v", vars, pRs.String())

	return fmt.Sprintf("%d", pRs.GetProcessInstanceKey()), nil

}

func (z *engineImpl) SendMessage(messageId string, correlationId string, vars map[string]interface{}) error {

	l := z.l().Mth("send-message").F(log.FF{"msg-id": messageId, "corr-id": correlationId})

	ctx := context.Background()

	m := z.client.NewPublishMessageCommand().MessageName(messageId).CorrelationKey(correlationId)

	if len(vars) > 0 {
		m, _ = m.VariablesFromMap(vars)
	}

	rs, err := m.Send(ctx)
	if err != nil {
		return ErrZeebeSendMessage(err, messageId, correlationId)
	}

	l.Dbg("sent").TrcF("vars=%v response=%v", vars, rs)

	return nil
}

func (z *engineImpl) SendError(jobId int64, errCode, errMessage string) error {

	l := z.l().Mth("send-err").F(log.FF{"jobId": jobId, "err-code": errCode, "err-m": errMessage})

	m := z.client.NewThrowErrorCommand().JobKey(jobId).ErrorCode(errCode).ErrorMessage(errMessage)

	_, err := m.Send(context.Background())
	if err != nil {
		return ErrZeebeSendError(err, jobId, errCode, errMessage)
	}
	l.Dbg("sent")

	return nil
}
