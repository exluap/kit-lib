package cron

import "github.com/exluap/kit/er"

const (
	ErrCodeCronStart = "CRON-001"
)

var (
	ErrCronStart = func(cause error, name string) error {
		return er.WrapWithBuilder(cause, ErrCodeCronStart, "").F(er.FF{"name": name}).Err()
	}
)
