package commentgrpc

import (
	"testing"
	"time"

	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	commentv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/comment"
	"github.com/stretchr/testify/assert"
)

func TestCommentToProto(t *testing.T) {
	comment := domain.Comment{
		ID:        "123",
		PostID:    "456",
		Body:      "Test comment",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		User:      domain.User{},
	}

	protoComment := CommentToProto(comment)

	assert.Equal(t, comment.ID, protoComment.Id)
	assert.Equal(t, comment.PostID, protoComment.PostId)
	assert.Equal(t, comment.Body, protoComment.Body)
	assert.Equal(t, comment.CreatedAt.String(), protoComment.CreatedAt)
	assert.Equal(t, comment.UpdatedAt.String(), protoComment.UpdatedAt)
}

func TestUserToProto(t *testing.T) {
	user := domain.User{
		ID:        123,
		FirstName: "John",
		LastName:  "Doe",
		AvatarURL: "http://example.com/avatar.jpg",
	}

	protoUser := UserToProto(user)

	assert.Equal(t, user.ID, protoUser.Id)
	assert.Equal(t, user.FirstName, protoUser.FirstName)
	assert.Equal(t, user.LastName, protoUser.LastName)
	assert.Equal(t, user.AvatarURL, protoUser.AvatarUrl)
}

func TestPaginationMetadataToProto(t *testing.T) {
	pagination := domain.PaginationMetadata{
		CurrentPage:  1,
		PageSize:     10,
		FirstPage:    1,
		LastPage:     10,
		TotalRecords: 100,
	}

	protoPagination := PaginationMetadataToProto(pagination)

	assert.Equal(t, pagination.CurrentPage, protoPagination.CurrentPage)
	assert.Equal(t, pagination.PageSize, protoPagination.PageSize)
	assert.Equal(t, pagination.FirstPage, protoPagination.FirstPage)
	assert.Equal(t, pagination.LastPage, protoPagination.LastPage)
	assert.Equal(t, pagination.TotalRecords, protoPagination.TotalRecords)
}

func TestCommentsToProto(t *testing.T) {
	comments := []domain.Comment{
		{
			ID:        "123",
			PostID:    "456",
			Body:      "Test comment 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			User:      domain.User{},
		},
		{
			ID:        "789",
			PostID:    "012",
			Body:      "Test comment 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			User:      domain.User{},
		},
	}

	protoComments := CommentsToProto(comments)

	assert.Equal(t, len(comments), len(protoComments))
	for i, comment := range comments {
		assert.Equal(t, comment.ID, protoComments[i].Id)
		assert.Equal(t, comment.PostID, protoComments[i].PostId)
		assert.Equal(t, comment.Body, protoComments[i].Body)
		assert.Equal(t, comment.CreatedAt.String(), protoComments[i].CreatedAt)
		assert.Equal(t, comment.UpdatedAt.String(), protoComments[i].UpdatedAt)
	}
}

func TestProtoToFilter(t *testing.T) {
	tests := []struct {
		name           string
		request        *commentv1.ListPostCommentsRequest
		expectedFilter domain.Filter
		expectedError  error
	}{
		{
			name: "Success",
			request: &commentv1.ListPostCommentsRequest{
				Page:      1,
				PageSize:  10,
				SortBy:    commentv1.SortBy_SORT_BY_CREATED_AT,
				SortOrder: commentv1.SortOrder_SORT_ORDER_DESC,
			},
			expectedFilter: domain.Filter{
				Page:      1,
				PageSize:  10,
				SortBy:    domain.SortByCreatedAt,
				SortOrder: domain.SortOrderDesc,
				FilterMap: map[string]bool{"page": true, "page_size": true, "sort_by": true, "sort_order": true},
			},
			expectedError: nil,
		},
		{
			name: "unspecified proto sort by must covert to unspecified domain",
			request: &commentv1.ListPostCommentsRequest{
				Page:      1,
				PageSize:  10,
				SortBy:    commentv1.SortBy_SORT_BY_UNSPECIFIED,
				SortOrder: commentv1.SortOrder_SORT_ORDER_ASC,
			},
			expectedFilter: domain.Filter{
				Page:      1,
				PageSize:  10,
				SortBy:    domain.SortByUnspecified,
				SortOrder: domain.SortOrderAsc,
				FilterMap: map[string]bool{"page": true, "page_size": true, "sort_by": true, "sort_order": true},
			},
			expectedError: nil,
		},
		{
			name: "unspecified proto sort order must covert to descending domain",
			request: &commentv1.ListPostCommentsRequest{
				Page:      1,
				PageSize:  10,
				SortBy:    commentv1.SortBy_SORT_BY_UPDATED_AT,
				SortOrder: commentv1.SortOrder_SORT_ORDER_UNSPECIFIED,
			},
			expectedFilter: domain.Filter{
				Page:      1,
				PageSize:  10,
				SortBy:    domain.SortByUpdatedAt,
				SortOrder: domain.SortOrderDesc,
				FilterMap: map[string]bool{"page": true, "page_size": true, "sort_by": true, "sort_order": true},
			},
			expectedError: nil,
		},
		{
			name: "empty page and page size must return default value: page=1, page_size=25",
			request: &commentv1.ListPostCommentsRequest{
				Page:      0,
				PageSize:  0,
				SortBy:    commentv1.SortBy_SORT_BY_CREATED_AT,
				SortOrder: commentv1.SortOrder_SORT_ORDER_DESC,
			},
			expectedFilter: domain.Filter{
				Page:      1,
				PageSize:  25,
				SortBy:    domain.SortByCreatedAt,
				SortOrder: domain.SortOrderDesc,
				FilterMap: map[string]bool{"page": true, "page_size": true, "sort_by": true, "sort_order": true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := ProtoToFilter(tt.request)

			assert.Equal(t, tt.expectedFilter, filter)
			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}
