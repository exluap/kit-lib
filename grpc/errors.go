package grpc

import (
	"context"
	"github.com/exluap/kit/er"
)

const (
	ErrCodeGrpcClientDial  = "GRPC-001"
	ErrCodeGrpCInvoke      = "GRPC-002"
	ErrCodeGrpcSrvListen   = "GRPC-003"
	ErrCodeGrpcSrvServe    = "GRPC-004"
	ErrCodeGrpcSrvNotReady = "GRPC-005"
	ErrCodeGrpcPanic       = "GRPC-006"
)

var (
	ErrGrpcClientDial  = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeGrpcClientDial, "").Err() }
	ErrGrpCInvoke      = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeGrpCInvoke, "").Err() }
	ErrGrpcSrvListen   = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeGrpcSrvListen, "").Err() }
	ErrGrpcSrvServe    = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeGrpcSrvServe, "").Err() }
	ErrGrpcSrvNotReady = func(svc string) error {
		return er.WithBuilder(ErrCodeGrpcSrvNotReady, "service isn't ready within timeout").F(er.FF{"svc": svc}).Err()
	}
	ErrGrpcPanic = func(ctx context.Context, cause interface{}) error {
		return er.WithBuilder(ErrCodeGrpcPanic, "panic").F(er.FF{"cause": cause}).C(ctx).Err()
	}
	ErrGrpcPanicNoCtx = func(cause interface{}) error {
		return er.WithBuilder(ErrCodeGrpcPanic, "panic").F(er.FF{"cause": cause}).Err()
	}
)
