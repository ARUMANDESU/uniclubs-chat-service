// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// Updater is an autogenerated mock type for the Updater type
type Updater struct {
	mock.Mock
}

// UpdateComment provides a mock function with given fields: ctx, comment
func (_m *Updater) UpdateComment(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
	ret := _m.Called(ctx, comment)

	if len(ret) == 0 {
		panic("no return value specified for UpdateComment")
	}

	var r0 domain.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Comment) (domain.Comment, error)); ok {
		return rf(ctx, comment)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.Comment) domain.Comment); ok {
		r0 = rf(ctx, comment)
	} else {
		r0 = ret.Get(0).(domain.Comment)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.Comment) error); ok {
		r1 = rf(ctx, comment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUpdater creates a new instance of Updater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUpdater(t interface {
	mock.TestingT
	Cleanup(func())
}) *Updater {
	mock := &Updater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
