package zeebe

import "github.com/exluap/kit/er"

const (
	ErrCodeZeebeSend               = "ZB-001"
	ErrCodeZeebeClose              = "ZB-002"
	ErrCodeZeebeConnInvalidParams  = "ZB-003"
	ErrCodeZeebeNewClient          = "ZB-004"
	ErrCodeZeebeConnInvalidHandler = "ZB-005"
	ErrCodeZeebeStartProcess       = "ZB-006"
	ErrCodeZeebeSendMessage        = "ZB-007"
	ErrCodeZeebeSendError          = "ZB-008"
	ErrCodeZeebeVarsFromMap        = "ZB-009"
	ErrCodeZeebeVarsAsMap          = "ZB-010"
	ErrCodeZeebeCtxInvalid         = "ZB-011"
	ErrCodeZeebeCtxNotFound        = "ZB-012"
)

var (
	ErrZeebeClose             = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeZeebeClose, "").Err() }
	ErrZeebeConnInvalidParams = func(host, port string) error {
		return er.WithBuilder(ErrCodeZeebeConnInvalidParams, "invalid params").F(er.FF{"host": host, "port": port}).Err()
	}
	ErrZeebeConnInvalidHandler = func() error { return er.WithBuilder(ErrCodeZeebeConnInvalidHandler, "invalid handler").Err() }
	ErrZeebeCtxInvalid         = func() error { return er.WithBuilder(ErrCodeZeebeCtxInvalid, "invalid _ctx vars").Err() }
	ErrZeebeCtxNotFound        = func() error { return er.WithBuilder(ErrCodeZeebeCtxNotFound, "not found _ctx vars").Err() }
	ErrZeebeNewClient          = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeZeebeNewClient, "").Err() }
	ErrZeebeStartProcess       = func(cause error, processId string) error {
		return er.WrapWithBuilder(cause, ErrCodeZeebeStartProcess, "").F(er.FF{"pid": processId}).Err()
	}
	ErrZeebeSendMessage = func(cause error, msgId, corrId string) error {
		return er.WrapWithBuilder(cause, ErrCodeZeebeSendMessage, "").F(er.FF{"msgId": msgId, "corrId": corrId}).Err()
	}
	ErrZeebeSendError = func(cause error, jobId int64, errCode, errMessage string) error {
		return er.WrapWithBuilder(cause, ErrCodeZeebeSendError, "").F(er.FF{"jobId": jobId, "errCode": errCode, "errMsg": errMessage}).Err()
	}
	ErrZeebeVarsFromMap = func(cause error, jobKey int64) error {
		return er.WrapWithBuilder(cause, ErrCodeZeebeVarsFromMap, "").F(er.FF{"jobKey": jobKey}).Err()
	}
	ErrZeebeVarsAsMap = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeZeebeVarsAsMap, "").Err() }
	ErrZeebeSend      = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeZeebeSend, "").Err() }
)
