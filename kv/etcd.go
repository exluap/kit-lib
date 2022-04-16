package kv

import (
	"github.com/exluap/kit/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"time"
)

const (
	Dial = time.Second * 3
)

type Etcd struct {
	Client *clientv3.Client
	logger log.CLoggerFunc
}

// Config etcd configs
type Config struct {
	Hosts []string
}

func Open(cfg *Config, logger log.CLoggerFunc) (*Etcd, error) {

	etcd := &Etcd{
		logger: logger,
	}

	cl, err := clientv3.New(clientv3.Config{
		DialTimeout: Dial,
		DialOptions: []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()},
		Endpoints:   cfg.Hosts,
	})
	if err != nil {
		return nil, ErrEtcdOpen(err)
	}

	logger().Cmp("etcd").Inf("ok")

	etcd.Client = cl
	return etcd, nil

}

func (e *Etcd) Close() error {
	e.logger().Cmp("etcd").Inf("closed")
	if e.Client != nil {
		return e.Client.Close()
	}
	return nil
}
