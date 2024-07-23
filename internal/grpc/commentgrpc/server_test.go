package commentgrpc

import (
	"context"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/grpc/commentgrpc/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	commentv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/comment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Suite struct {
	mockService *mocks.CommentService
	server      Server
}

func NewSuite(t *testing.T) *Suite {
	mockService := mocks.NewCommentService(t)
	server := NewServer(mockService)
	return &Suite{mockService: mockService, server: server}
}

func TestServer_GetCommentByID(t *testing.T) {
	st := NewSuite(t)
	request := &commentv1.GetCommentByIDRequest{Id: "123"}

	st.mockService.On("GetByID", mock.Anything, "123").Return(domain.Comment{ID: "123"}, nil)

	response, err := st.server.GetCommentByID(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, "123", response.Comment.Id)
	st.mockService.AssertExpectations(t)
}

func TestServer_GetCommentByID_InvalidArgument(t *testing.T) {
	st := NewSuite(t)
	request := &commentv1.GetCommentByIDRequest{}

	response, err := st.server.GetCommentByID(context.Background(), request)

	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Nil(t, response)
}

func TestServer_ListPostComments(t *testing.T) {
	st := NewSuite(t)
	request := &commentv1.ListPostCommentsRequest{PostId: "123", Page: 1, PageSize: 10}

	st.mockService.On("ListByPostID", mock.Anything, "123", mock.Anything).Return([]domain.Comment{{ID: "123"}}, domain.PaginationMetadata{}, nil)

	response, err := st.server.ListPostComments(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, "123", response.Comments[0].Id)
	st.mockService.AssertExpectations(t)
}
