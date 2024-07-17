package userservice

import (
	"context"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/services/userservice/mocks"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type Suite struct {
	Service               Service
	mockPrimaryProvider   *mocks.UserProvider
	mockSecondaryProvider *mocks.UserProvider
	mockSaver             *mocks.UserSaver
}

func NewSuite(t *testing.T) *Suite {
	s := &Suite{
		mockPrimaryProvider:   mocks.NewUserProvider(t),
		mockSecondaryProvider: mocks.NewUserProvider(t),
		mockSaver:             mocks.NewUserSaver(t),
	}
	s.Service = New(logger.Plug(), s.mockPrimaryProvider, s.mockSecondaryProvider, s.mockSaver)
	return s
}

func TestService_GetUser_PrimaryProviderSuccess(t *testing.T) {
	s := NewSuite(t)
	defer s.mockPrimaryProvider.AssertExpectations(t)

	s.mockPrimaryProvider.On("GetUserByID", mock.Anything, int64(1)).Return(domain.User{}, nil)

	user, err := s.Service.GetUser(context.Background(), int64(1))
	assert.Nil(t, err)
	assert.NotNil(t, user)

}

func TestService_GetUser_PrimaryProviderFail(t *testing.T) {
	s := NewSuite(t)
	defer s.mockPrimaryProvider.AssertExpectations(t)
	defer s.mockSecondaryProvider.AssertExpectations(t)
	defer s.mockSaver.AssertExpectations(t)

	s.mockPrimaryProvider.On("GetUserByID", mock.Anything, int64(1)).Return(domain.User{}, ErrUserNotFound)
	s.mockSecondaryProvider.On("GetUserByID", mock.Anything, int64(1)).Return(domain.User{}, nil)
	s.mockSaver.On("SaveUser", mock.Anything, mock.Anything).Return(nil)

	user, err := s.Service.GetUser(context.Background(), int64(1))
	assert.Nil(t, err)
	assert.NotNil(t, user)
}

func TestService_GetUser_AllProvidersFail(t *testing.T) {
	s := NewSuite(t)
	defer s.mockPrimaryProvider.AssertExpectations(t)
	defer s.mockSecondaryProvider.AssertExpectations(t)
	defer s.mockSaver.AssertExpectations(t)

	s.mockPrimaryProvider.On("GetUserByID", mock.Anything, int64(1)).Return(domain.User{}, ErrUserNotFound)
	s.mockSecondaryProvider.On("GetUserByID", mock.Anything, int64(1)).Return(domain.User{}, ErrUserNotFound)

	_, err := s.Service.GetUser(context.Background(), int64(1))
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrUserNotFound)
}
