package commentservice

import (
	"context"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/services/commentservice/mocks"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type Suite struct {
	Service          Service
	mockProvider     *mocks.Provider
	mockCreator      *mocks.Creator
	mockUpdater      *mocks.Updater
	mockDeleter      *mocks.Deleter
	mockUserProvider *mocks.UserProvider
}

func newSuite(t *testing.T) *Suite {
	s := &Suite{
		mockProvider:     mocks.NewProvider(t),
		mockCreator:      mocks.NewCreator(t),
		mockUpdater:      mocks.NewUpdater(t),
		mockDeleter:      mocks.NewDeleter(t),
		mockUserProvider: mocks.NewUserProvider(t),
	}
	s.Service = New(Config{
		Logger:       logger.Plug(),
		Provider:     s.mockProvider,
		Creator:      s.mockCreator,
		Updater:      s.mockUpdater,
		Deleter:      s.mockDeleter,
		UserProvider: s.mockUserProvider,
	})
	return s
}

func TestService_Create(t *testing.T) {
	s := newSuite(t)
	defer s.mockCreator.AssertExpectations(t)
	defer s.mockUserProvider.AssertExpectations(t)

	s.mockUserProvider.On("GetUser", mock.Anything, mock.Anything).Return(domain.User{}, nil)
	s.mockCreator.On("CreateComment", mock.Anything, mock.Anything).Return(domain.Comment{}, nil)

	comment, err := s.Service.Create(context.Background(), CreateCommentDTO{})
	assert.Nil(t, err)
	assert.NotNil(t, comment)
}

func TestService_Create_FailPath(t *testing.T) {

	testCases := []struct {
		name          string
		onGetUser     error
		onCreate      error
		expectedError error
	}{
		{
			name:          "unexpected error",
			onGetUser:     assert.AnError,
			onCreate:      nil,
			expectedError: domain.ErrInternal,
		},
		{
			name:          "user not found",
			onGetUser:     domain.ErrUserNotFound,
			onCreate:      assert.AnError,
			expectedError: domain.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := newSuite(t)

			defer s.mockUserProvider.AssertExpectations(t)
			s.mockUserProvider.On("GetUser", mock.Anything, mock.Anything).Return(domain.User{}, tc.onGetUser)

			if tc.onGetUser == nil {
				defer s.mockCreator.AssertExpectations(t)
				s.mockCreator.On("CreateComment", mock.Anything, mock.Anything).Return(domain.Comment{}, tc.onCreate)
			}

			_, err := s.Service.Create(context.Background(), CreateCommentDTO{})

			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
