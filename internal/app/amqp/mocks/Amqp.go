// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	amqp091 "github.com/rabbitmq/amqp091-go"

	mock "github.com/stretchr/testify/mock"
)

// Amqp is an autogenerated mock type for the Amqp type
type Amqp struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Amqp) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Consume provides a mock function with given fields: queue, routingKey, handler
func (_m *Amqp) Consume(queue string, routingKey string, handler func(amqp091.Delivery) error) error {
	ret := _m.Called(queue, routingKey, handler)

	if len(ret) == 0 {
		panic("no return value specified for Consume")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, func(amqp091.Delivery) error) error); ok {
		r0 = rf(queue, routingKey, handler)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewAmqp creates a new instance of Amqp. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAmqp(t interface {
	mock.TestingT
	Cleanup(func())
}) *Amqp {
	mock := &Amqp{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
