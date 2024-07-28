package userservice

import (
	"context"
	"errors"
	"testing"

	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/services/userservice/mocks"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Suite struct {
	Service      Service
	dbProvider   *mocks.UserProvider
	gRPCProvider *mocks.UserGRPCProvider
}

func NewSuite(t *testing.T) *Suite {
	s := &Suite{
		dbProvider:   mocks.NewUserProvider(t),
		gRPCProvider: mocks.NewUserGRPCProvider(t),
	}
	s.Service = New(logger.Plug(), s.dbProvider, s.gRPCProvider)
	return s
}

func TestService_GetUser(t *testing.T) {
	tests := []struct {
		name                   string
		userId                 int64
		primaryProviderError   error
		secondaryProviderError error
		expectedError          error
	}{
		{
			name:                   "Primary provider success",
			userId:                 1,
			primaryProviderError:   nil,
			secondaryProviderError: nil,
			expectedError:          nil,
		},
		{
			name:                   "Primary provider fail",
			userId:                 1,
			primaryProviderError:   domain.ErrUserNotFound,
			secondaryProviderError: nil,
			expectedError:          nil,
		},
		{
			name:                   "All providers fail",
			userId:                 1,
			primaryProviderError:   domain.ErrUserNotFound,
			secondaryProviderError: domain.ErrUserNotFound,
			expectedError:          domain.ErrUserNotFound,
		},
		{
			name:                   "Unexpected error, but secondary provider success",
			userId:                 1,
			primaryProviderError:   errors.New("unexpected error"),
			secondaryProviderError: nil,
			expectedError:          nil,
		},
		{
			name:                   "Unexpected error",
			userId:                 1,
			primaryProviderError:   errors.New("unexpected error"),
			secondaryProviderError: errors.New("unexpected error"),
			expectedError:          domain.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)

			defer s.dbProvider.AssertExpectations(t)
			s.dbProvider.On("GetUserByID", mock.Anything, tt.userId).Return(domain.User{}, tt.primaryProviderError)

			if tt.primaryProviderError != nil {
				defer s.gRPCProvider.AssertExpectations(t)
				s.gRPCProvider.On("GetUserByID", mock.Anything, tt.userId).Return(domain.User{}, tt.secondaryProviderError)
			}

			_, err := s.Service.GetUser(context.Background(), tt.userId)
			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name          string
		user          domain.User
		primaryError  error
		expectedError error
	}{
		{
			name:          "success",
			user:          domain.User{ID: 1},
			primaryError:  nil,
			expectedError: nil,
		},
		{
			name:          "fail: not found",
			user:          domain.User{ID: 1},
			primaryError:  domain.ErrUserNotFound,
			expectedError: domain.ErrUserNotFound,
		},
		{
			name:          "unexpected error",
			user:          domain.User{ID: 1},
			primaryError:  errors.New("unexpected error"),
			expectedError: domain.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)

			s.dbProvider.On("UpdateUser", mock.Anything, tt.user).Return(tt.primaryError)
			defer s.dbProvider.AssertExpectations(t)

			err := s.Service.Update(context.Background(), tt.user)
			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}
