// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// Provider is an autogenerated mock type for the Provider type
type Provider struct {
	mock.Mock
}

// GetComment provides a mock function with given fields: ctx, commentID
func (_m *Provider) GetComment(ctx context.Context, commentID string) (domain.Comment, error) {
	ret := _m.Called(ctx, commentID)

	if len(ret) == 0 {
		panic("no return value specified for GetComment")
	}

	var r0 domain.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.Comment, error)); ok {
		return rf(ctx, commentID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Comment); ok {
		r0 = rf(ctx, commentID)
	} else {
		r0 = ret.Get(0).(domain.Comment)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, commentID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPostComments provides a mock function with given fields: ctx, postID, filter
func (_m *Provider) ListPostComments(ctx context.Context, postID string, filter domain.Filter) ([]domain.Comment, domain.PaginationMetadata, error) {
	ret := _m.Called(ctx, postID, filter)

	if len(ret) == 0 {
		panic("no return value specified for ListPostComments")
	}

	var r0 []domain.Comment
	var r1 domain.PaginationMetadata
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.Filter) ([]domain.Comment, domain.PaginationMetadata, error)); ok {
		return rf(ctx, postID, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.Filter) []domain.Comment); ok {
		r0 = rf(ctx, postID, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Comment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, domain.Filter) domain.PaginationMetadata); ok {
		r1 = rf(ctx, postID, filter)
	} else {
		r1 = ret.Get(1).(domain.PaginationMetadata)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, domain.Filter) error); ok {
		r2 = rf(ctx, postID, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewProvider creates a new instance of Provider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *Provider {
	mock := &Provider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
