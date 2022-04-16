package kv

import "github.com/exluap/kit/er"

const (
	ErrCodeEtcdOpen = "ETCD-001"
)

var (
	ErrEtcdOpen = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeEtcdOpen, "").Err() }
)
