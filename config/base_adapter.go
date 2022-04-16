package config

import (
	kitGrpc "github.com/exluap/kit/grpc"
	kitLog "github.com/exluap/kit/log"
	"time"
)

type BaseAdapter struct {
	Client *kitGrpc.Client
	LogFn  kitLog.CLoggerFunc
}

func (a *BaseAdapter) AwaitReadiness(timeout time.Duration) error {
	l := a.LogFn().Mth("await").DbgF("awaiting config-server readiness, timeout=%v", timeout)
	if !a.Client.AwaitReadiness(timeout) {
		return ErrConfigTimeout()
	}
	l.Dbg("ready")
	return nil
}

func (a *BaseAdapter) Close() {
	_ = a.Client.Conn.Close()
}
