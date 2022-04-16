package redis

import "github.com/exluap/kit/er"

const (
	ErrCodeRedisPingErr = "RDS-001"
)

var (
	ErrRedisPingErr = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeRedisPingErr, "").Err() }
)
