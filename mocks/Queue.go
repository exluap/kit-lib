// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	queue "github.com/exluap/kit/queue"
	mock "github.com/stretchr/testify/mock"
)

// Queue is an autogenerated mock type for the Queue type
type Queue struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Queue) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Open provides a mock function with given fields: ctx, clientId, options
func (_m *Queue) Open(ctx context.Context, clientId string, options *queue.Config) error {
	ret := _m.Called(ctx, clientId, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *queue.Config) error); ok {
		r0 = rf(ctx, clientId, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Publish provides a mock function with given fields: ctx, qt, topic, msg
func (_m *Queue) Publish(ctx context.Context, qt queue.QueueType, topic string, msg *queue.Message) error {
	ret := _m.Called(ctx, qt, topic, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, queue.QueueType, string, *queue.Message) error); ok {
		r0 = rf(ctx, qt, topic, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Subscribe provides a mock function with given fields: qt, topic, receiverChan
func (_m *Queue) Subscribe(qt queue.QueueType, topic string, receiverChan chan<- []byte) error {
	ret := _m.Called(qt, topic, receiverChan)

	var r0 error
	if rf, ok := ret.Get(0).(func(queue.QueueType, string, chan<- []byte) error); ok {
		r0 = rf(qt, topic, receiverChan)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubscribeLB provides a mock function with given fields: qt, topic, loadBalancingGroup, receiverChan
func (_m *Queue) SubscribeLB(qt queue.QueueType, topic string, loadBalancingGroup string, receiverChan chan<- []byte) error {
	ret := _m.Called(qt, topic, loadBalancingGroup, receiverChan)

	var r0 error
	if rf, ok := ret.Get(0).(func(queue.QueueType, string, string, chan<- []byte) error); ok {
		r0 = rf(qt, topic, loadBalancingGroup, receiverChan)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
