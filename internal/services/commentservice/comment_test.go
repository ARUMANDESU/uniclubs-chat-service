package commentservice

import (
	"context"
	"testing"
	"time"

	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/services/commentservice/mocks"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
				UserID:    1,
				Body:      "new body",
			},
		},
		{
			name: "empty comment",
			dto: UpdateCommentDTO{
				UserID:    1,
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
		dto           UpdateCommentDTO
		onGetComment  error
		onUpdate      error
		expectedError error
	}{
		{
			name:          "unexpected error",
			dto:           UpdateCommentDTO{CommentID: "1", UserID: 1, Body: "new body"},
			onGetComment:  nil,
			onUpdate:      assert.AnError,
			expectedError: domain.ErrInternal,
		},
		{
			name:          "comment not found",
			dto:           UpdateCommentDTO{CommentID: "1", UserID: 1, Body: "new body"},
			onGetComment:  domain.ErrCommentNotFound,
			onUpdate:      nil,
			expectedError: domain.ErrCommentNotFound,
		},
		{
			name:          "invalid id",
			dto:           UpdateCommentDTO{CommentID: "1", UserID: 1, Body: "new body"},
			onGetComment:  domain.ErrInvalidID,
			onUpdate:      nil,
			expectedError: domain.ErrInvalidID,
		},
		{
			name:          "unauthorized",
			dto:           UpdateCommentDTO{CommentID: "1", UserID: 2, Body: "new body"},
			onGetComment:  nil,
			onUpdate:      nil,
			expectedError: domain.ErrUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newSuite(t)
			defer s.mockProvider.AssertExpectations(t)
			defer s.mockUpdater.AssertExpectations(t)

			s.mockProvider.On("GetComment", mock.Anything, "1").Return(baseComment, tc.onGetComment)

			if tc.onGetComment == nil && tc.expectedError != domain.ErrUnauthorized {
				s.mockUpdater.On("UpdateComment", mock.Anything, mock.AnythingOfType("domain.Comment")).Return(func(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
					return comment, tc.onUpdate
				})
			}

			_, err := s.Service.Update(context.Background(), tc.dto)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestService_Delete(t *testing.T) {
	s := newSuite(t)
	defer s.mockProvider.AssertExpectations(t)
	defer s.mockDeleter.AssertExpectations(t)

	s.mockProvider.On("GetComment", mock.Anything, "1").Return(domain.Comment{User: domain.User{ID: 1}}, nil)
	s.mockDeleter.On("DeleteComment", mock.Anything, "1").Return(nil)

	err := s.Service.Delete(context.Background(), DeleteCommentDTO{
		UserID:    1,
		CommentID: "1",
	})
	assert.Nil(t, err)
}

func TestService_Delete_FailPath(t *testing.T) {

	tests := []struct {
		name          string
		dto           DeleteCommentDTO
		onGetComment  error
		onDelete      error
		expectedError error
	}{
		{
			name:          "unexpected error",
			dto:           DeleteCommentDTO{UserID: 0, CommentID: "1"},
			onGetComment:  nil,
			onDelete:      assert.AnError,
			expectedError: domain.ErrInternal,
		},
		{
			name:          "comment not found",
			dto:           DeleteCommentDTO{UserID: 0, CommentID: "1"},
			onGetComment:  domain.ErrCommentNotFound,
			onDelete:      nil,
			expectedError: domain.ErrCommentNotFound,
		},
		{
			name:          "invalid id",
			dto:           DeleteCommentDTO{UserID: 0, CommentID: "1"},
			onGetComment:  domain.ErrInvalidID,
			onDelete:      nil,
			expectedError: domain.ErrInvalidID,
		},
		{
			name:          "unauthorized",
			dto:           DeleteCommentDTO{UserID: 1, CommentID: "1"},
			onGetComment:  nil,
			onDelete:      nil,
			expectedError: domain.ErrUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newSuite(t)
			defer s.mockProvider.AssertExpectations(t)
			defer s.mockDeleter.AssertExpectations(t)

			s.mockProvider.On("GetComment", mock.Anything, "1").Return(domain.Comment{}, tc.onGetComment)

			if tc.onGetComment == nil && tc.expectedError != domain.ErrUnauthorized {
				s.mockDeleter.On("DeleteComment", mock.Anything, "1").Return(tc.onDelete)
			}

			err := s.Service.Delete(context.Background(), tc.dto)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestService_GetByID(t *testing.T) {
	s := newSuite(t)
	defer s.mockProvider.AssertExpectations(t)

	s.mockProvider.On("GetComment", mock.Anything, "1").Return(domain.Comment{}, nil)

	comment, err := s.Service.GetByID(context.Background(), "1")
	assert.Nil(t, err)
	assert.NotNil(t, comment)
}

func TestService_GetByID_FailPath(t *testing.T) {

	tests := []struct {
		name          string
		onGetComment  error
		expectedError error
	}{
		{
			name:          "unexpected error",
			onGetComment:  assert.AnError,
			expectedError: domain.ErrInternal,
		},
		{
			name:          "comment not found",
			onGetComment:  domain.ErrCommentNotFound,
			expectedError: domain.ErrCommentNotFound,
		},
		{
			name:          "invalid id",
			onGetComment:  domain.ErrInvalidID,
			expectedError: domain.ErrInvalidID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newSuite(t)
			defer s.mockProvider.AssertExpectations(t)

			s.mockProvider.On("GetComment", mock.Anything, "1").Return(domain.Comment{}, tc.onGetComment)

			_, err := s.Service.GetByID(context.Background(), "1")
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestService_ListByPostID(t *testing.T) {
	s := newSuite(t)
	defer s.mockProvider.AssertExpectations(t)

	filter, err := domain.NewFilter(domain.WithPage(1), domain.WithPageSize(10))
	assert.Nil(t, err)

	expectedComments := []domain.Comment{
		{ID: "1", Body: "Comment 1"},
		{ID: "2", Body: "Comment 2"},
	}
	expectedMetadata := domain.PaginationMetadata{TotalRecords: 2, PageSize: 10, CurrentPage: 1, FirstPage: 1, LastPage: 1}

	s.mockProvider.On("ListPostComments", mock.Anything, "1", mock.Anything).Return(expectedComments, expectedMetadata, nil)

	comments, metadata, err := s.Service.ListByPostID(context.Background(), "1", *filter)
	assert.Nil(t, err)
	assert.Equal(t, expectedComments, comments)
	assert.Equal(t, expectedMetadata, metadata)
}

func TestService_ListByPostID_FailPath(t *testing.T) {
	tests := []struct {
		name               string
		onListPostComments error
		expectedError      error
	}{
		{
			name:               "unexpected error",
			onListPostComments: assert.AnError,
			expectedError:      domain.ErrInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newSuite(t)
			defer s.mockProvider.AssertExpectations(t)

			s.mockProvider.On("ListPostComments", mock.Anything, "1", domain.Filter{}).Return(nil, domain.PaginationMetadata{}, tc.onListPostComments)

			_, _, err := s.Service.ListByPostID(context.Background(), "1", domain.Filter{})
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
