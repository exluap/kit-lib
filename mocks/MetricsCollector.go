// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	monitoring "github.com/exluap/kit/monitoring"
	mock "github.com/stretchr/testify/mock"
)

// MetricsCollector is an autogenerated mock type for the MetricsCollector type
type MetricsCollector struct {
	mock.Mock
}

// Execute provides a mock function with given fields:
func (_m *MetricsCollector) Execute() monitoring.MetricsCollection {
	ret := _m.Called()

	var r0 monitoring.MetricsCollection
	if rf, ok := ret.Get(0).(func() monitoring.MetricsCollection); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(monitoring.MetricsCollection)
		}
	}

	return r0
}
