// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	log "github.com/exluap/kit/log"
	mock "github.com/stretchr/testify/mock"
)

// CLogger is an autogenerated mock type for the CLogger type
type CLogger struct {
	mock.Mock
}

// C provides a mock function with given fields: ctx
func (_m *CLogger) C(ctx context.Context) log.CLogger {
	ret := _m.Called(ctx)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(context.Context) log.CLogger); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// Clone provides a mock function with given fields:
func (_m *CLogger) Clone() log.CLogger {
	ret := _m.Called()

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func() log.CLogger); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// Cmp provides a mock function with given fields: c
func (_m *CLogger) Cmp(c string) log.CLogger {
	ret := _m.Called(c)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(string) log.CLogger); ok {
		r0 = rf(c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// Dbg provides a mock function with given fields: args
func (_m *CLogger) Dbg(args ...interface{}) log.CLogger {
	var _ca []interface{}
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(...interface{}) log.CLogger); ok {
		r0 = rf(args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// DbgF provides a mock function with given fields: format, args
func (_m *CLogger) DbgF(format string, args ...interface{}) log.CLogger {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(string, ...interface{}) log.CLogger); ok {
		r0 = rf(format, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// E provides a mock function with given fields: err
func (_m *CLogger) E(err error) log.CLogger {
	ret := _m.Called(err)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(error) log.CLogger); ok {
		r0 = rf(err)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// Err provides a mock function with given fields: args
func (_m *CLogger) Err(args ...interface{}) log.CLogger {
	var _ca []interface{}
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(...interface{}) log.CLogger); ok {
		r0 = rf(args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// ErrF provides a mock function with given fields: format, args
func (_m *CLogger) ErrF(format string, args ...interface{}) log.CLogger {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(string, ...interface{}) log.CLogger); ok {
		r0 = rf(format, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// F provides a mock function with given fields: fields
func (_m *CLogger) F(fields log.FF) log.CLogger {
	ret := _m.Called(fields)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(log.FF) log.CLogger); ok {
		r0 = rf(fields)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// Fatal provides a mock function with given fields: args
func (_m *CLogger) Fatal(args ...interface{}) log.CLogger {
	var _ca []interface{}
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(...interface{}) log.CLogger); ok {
		r0 = rf(args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// FatalF provides a mock function with given fields: format, args
func (_m *CLogger) FatalF(format string, args ...interface{}) log.CLogger {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(string, ...interface{}) log.CLogger); ok {
		r0 = rf(format, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// Inf provides a mock function with given fields: args
func (_m *CLogger) Inf(args ...interface{}) log.CLogger {
	var _ca []interface{}
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(...interface{}) log.CLogger); ok {
		r0 = rf(args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// InfF provides a mock function with given fields: format, args
func (_m *CLogger) InfF(format string, args ...interface{}) log.CLogger {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(string, ...interface{}) log.CLogger); ok {
		r0 = rf(format, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// Mth provides a mock function with given fields: m
func (_m *CLogger) Mth(m string) log.CLogger {
	ret := _m.Called(m)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(string) log.CLogger); ok {
		r0 = rf(m)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// Node provides a mock function with given fields: n
func (_m *CLogger) Node(n string) log.CLogger {
	ret := _m.Called(n)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(string) log.CLogger); ok {
		r0 = rf(n)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// Pr provides a mock function with given fields: m
func (_m *CLogger) Pr(m string) log.CLogger {
	ret := _m.Called(m)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(string) log.CLogger); ok {
		r0 = rf(m)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// Printf provides a mock function with given fields: _a0, _a1
func (_m *CLogger) Printf(_a0 string, _a1 ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _a1...)
	_m.Called(_ca...)
}

// Srv provides a mock function with given fields: s
func (_m *CLogger) Srv(s string) log.CLogger {
	ret := _m.Called(s)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(string) log.CLogger); ok {
		r0 = rf(s)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// St provides a mock function with given fields:
func (_m *CLogger) St() log.CLogger {
	ret := _m.Called()

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func() log.CLogger); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// Trc provides a mock function with given fields: args
func (_m *CLogger) Trc(args ...interface{}) log.CLogger {
	var _ca []interface{}
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(...interface{}) log.CLogger); ok {
		r0 = rf(args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// TrcF provides a mock function with given fields: format, args
func (_m *CLogger) TrcF(format string, args ...interface{}) log.CLogger {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(string, ...interface{}) log.CLogger); ok {
		r0 = rf(format, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// TrcObj provides a mock function with given fields: format, args
func (_m *CLogger) TrcObj(format string, args ...interface{}) log.CLogger {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(string, ...interface{}) log.CLogger); ok {
		r0 = rf(format, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// Warn provides a mock function with given fields: args
func (_m *CLogger) Warn(args ...interface{}) log.CLogger {
	var _ca []interface{}
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(...interface{}) log.CLogger); ok {
		r0 = rf(args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// WarnF provides a mock function with given fields: format, args
func (_m *CLogger) WarnF(format string, args ...interface{}) log.CLogger {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 log.CLogger
	if rf, ok := ret.Get(0).(func(string, ...interface{}) log.CLogger); ok {
		r0 = rf(format, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.CLogger)
		}
	}

	return r0
}

// Write provides a mock function with given fields: p
func (_m *CLogger) Write(p []byte) (int, error) {
	ret := _m.Called(p)

	var r0 int
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
