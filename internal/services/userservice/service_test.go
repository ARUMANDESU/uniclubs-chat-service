package userservice

import (
	"context"
	"errors"
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
}

func NewSuite(t *testing.T) *Suite {
	s := &Suite{
		mockPrimaryProvider:   mocks.NewUserProvider(t),
		mockSecondaryProvider: mocks.NewUserProvider(t),
	}
	s.Service = New(logger.Plug(), s.mockPrimaryProvider, s.mockSecondaryProvider)
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

			defer s.mockPrimaryProvider.AssertExpectations(t)
			s.mockPrimaryProvider.On("GetUserByID", mock.Anything, tt.userId).Return(domain.User{}, tt.primaryProviderError)

			if tt.primaryProviderError != nil {
				defer s.mockSecondaryProvider.AssertExpectations(t)
				s.mockSecondaryProvider.On("GetUserByID", mock.Anything, tt.userId).Return(domain.User{}, tt.secondaryProviderError)
			}

			_, err := s.Service.GetUser(context.Background(), tt.userId)
			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}
