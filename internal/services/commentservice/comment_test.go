package commentservice

import (
	"context"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/services/commentservice/mocks"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
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

func TestService_Update(t *testing.T) {
	baseComment := domain.Comment{
		ID:     "1",
		PostID: "1",
		User: domain.User{
			ID: 1,
		},
		Body:      "body",
		CreatedAt: time.Now().AddDate(0, 0, -2),
		UpdatedAt: time.Now().AddDate(0, 0, -2),
	}

	tests := []struct {
		name string
		dto  UpdateCommentDTO
	}{
		{
			name: "success",
			dto: UpdateCommentDTO{
				CommentID: "1",
				Body:      "new body",
			},
		},
		{
			name: "empty comment",
			dto: UpdateCommentDTO{
				CommentID: "1",
				Body:      "",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newSuite(t)
			defer s.mockProvider.AssertExpectations(t)
			defer s.mockUpdater.AssertExpectations(t)

			s.mockProvider.On("GetComment", mock.Anything, "1").Return(baseComment, nil)
			s.mockUpdater.On("UpdateComment", mock.Anything, mock.AnythingOfType("domain.Comment")).Return(func(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
				return comment, nil
			})

			comment, err := s.Service.Update(context.Background(), tc.dto)
			assert.Nil(t, err)
			assert.NotNil(t, comment)
			assert.Equal(t, tc.dto.Body, comment.Body)
			assert.NotEqual(t, baseComment.UpdatedAt, comment.UpdatedAt)
		})
	}
}

func TestService_Update_FailPath(t *testing.T) {

	baseComment := domain.Comment{
		ID:     "1",
		PostID: "1",
		User: domain.User{
			ID: 1,
		},
		Body:      "body",
		CreatedAt: time.Now().AddDate(0, 0, -2),
		UpdatedAt: time.Now().AddDate(0, 0, -2),
	}

	tests := []struct {
		name          string
		onGetComment  error
		onUpdate      error
		expectedError error
	}{
		{
			name:          "unexpected error",
			onGetComment:  nil,
			onUpdate:      assert.AnError,
			expectedError: domain.ErrInternal,
		},
		{
			name:          "comment not found",
			onGetComment:  domain.ErrCommentNotFound,
			onUpdate:      nil,
			expectedError: domain.ErrCommentNotFound,
		},
		{
			name:          "invalid id",
			onGetComment:  domain.ErrInvalidID,
			onUpdate:      nil,
			expectedError: domain.ErrInvalidID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newSuite(t)
			defer s.mockProvider.AssertExpectations(t)
			defer s.mockUpdater.AssertExpectations(t)

			s.mockProvider.On("GetComment", mock.Anything, "1").Return(baseComment, tc.onGetComment)

			if tc.onGetComment == nil {
				s.mockUpdater.On("UpdateComment", mock.Anything, mock.AnythingOfType("domain.Comment")).Return(func(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
					return comment, tc.onUpdate
				})
			}

			_, err := s.Service.Update(context.Background(), UpdateCommentDTO{CommentID: "1"})
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestService_Delete(t *testing.T) {
	s := newSuite(t)
	defer s.mockProvider.AssertExpectations(t)
	defer s.mockDeleter.AssertExpectations(t)

	s.mockProvider.On("GetComment", mock.Anything, "1").Return(domain.Comment{}, nil)
	s.mockDeleter.On("DeleteComment", mock.Anything, "1").Return(nil)

	err := s.Service.Delete(context.Background(), "1")
	assert.Nil(t, err)
}

func TestService_Delete_FailPath(t *testing.T) {

	tests := []struct {
		name          string
		onGetComment  error
		onDelete      error
		expectedError error
	}{
		{
			name:          "unexpected error",
			onGetComment:  nil,
			onDelete:      assert.AnError,
			expectedError: domain.ErrInternal,
		},
		{
			name:          "comment not found",
			onGetComment:  domain.ErrCommentNotFound,
			onDelete:      nil,
			expectedError: domain.ErrCommentNotFound,
		},
		{
			name:          "invalid id",
			onGetComment:  domain.ErrInvalidID,
			onDelete:      nil,
			expectedError: domain.ErrInvalidID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newSuite(t)
			defer s.mockProvider.AssertExpectations(t)
			defer s.mockDeleter.AssertExpectations(t)

			s.mockProvider.On("GetComment", mock.Anything, "1").Return(domain.Comment{}, tc.onGetComment)

			if tc.onGetComment == nil {
				s.mockDeleter.On("DeleteComment", mock.Anything, "1").Return(tc.onDelete)
			}

			err := s.Service.Delete(context.Background(), "1")
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
