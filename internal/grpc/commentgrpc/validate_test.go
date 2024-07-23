package commentgrpc

import (
	commentv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/comment"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateListPostComments_NilRequest(t *testing.T) {
	err := validateListPostComments(nil)
	assert.Error(t, err)
	assert.Equal(t, "request is nil", err.Error())
}

func TestValidateListPostComments_EmptyPostId(t *testing.T) {
	request := &commentv1.ListPostCommentsRequest{}
	err := validateListPostComments(request)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "post id is required")
}

func TestValidateListPostComments_NegativePage(t *testing.T) {
	request := &commentv1.ListPostCommentsRequest{
		PostId: "123",
		Page:   -1,
	}
	err := validateListPostComments(request)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "page must be greater than or equal to 0")
}

func TestValidateListPostComments_NegativePageSize(t *testing.T) {
	request := &commentv1.ListPostCommentsRequest{
		PostId:   "123",
		Page:     1,
		PageSize: -1,
	}
	err := validateListPostComments(request)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "page size must be greater than or equal to 0")
}

func TestValidateListPostComments_InvalidSortBy(t *testing.T) {
	request := &commentv1.ListPostCommentsRequest{
		PostId:   "123",
		Page:     1,
		PageSize: 10,
		SortBy:   999, // Invalid SortBy value
	}
	err := validateListPostComments(request)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SortBy: value must be one of")
}

func TestValidateListPostComments_InvalidSortOrder(t *testing.T) {
	request := &commentv1.ListPostCommentsRequest{
		PostId:    "123",
		Page:      1,
		PageSize:  10,
		SortBy:    commentv1.SortBy_SORT_BY_CREATED_AT,
		SortOrder: 999, // Invalid SortOrder value
	}
	err := validateListPostComments(request)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SortOrder: value must be one of")
}

func TestValidateListPostComments_ValidRequest(t *testing.T) {
	request := &commentv1.ListPostCommentsRequest{
		PostId:    "123",
		Page:      1,
		PageSize:  10,
		SortBy:    commentv1.SortBy_SORT_BY_CREATED_AT,
		SortOrder: commentv1.SortOrder_SORT_ORDER_ASC,
	}
	err := validateListPostComments(request)
	assert.NoError(t, err)
}
