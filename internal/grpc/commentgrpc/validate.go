package commentgrpc

import (
	"errors"
	"fmt"
	commentv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/comment"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func validateListPostComments(request *commentv1.ListPostCommentsRequest) error {
	if request == nil {
		return errors.New("request is nil")
	}

	validSortBy := []any{commentv1.SortBy_SORT_BY_UNSPECIFIED, commentv1.SortBy_SORT_BY_CREATED_AT, commentv1.SortBy_SORT_BY_UPDATED_AT}
	validSortOrder := []any{commentv1.SortOrder_SORT_ORDER_UNSPECIFIED, commentv1.SortOrder_SORT_ORDER_ASC, commentv1.SortOrder_SORT_ORDER_DESC}

	err := validation.ValidateStruct(request,
		validation.Field(&request.PostId, validation.Required.Error("post id is required")),
		validation.Field(&request.Page, validation.Min(0).Error("page must be greater than or equal to 0")),
		validation.Field(&request.PageSize, validation.Min(0).Error("page size must be greater than or equal to 0")),
		validation.Field(
			&request.SortBy,
			validation.In(validSortBy...).Error(fmt.Sprintf("SortBy: value must be one of %v", validSortBy)),
		),
		validation.Field(
			&request.SortOrder,
			validation.In(validSortOrder...).Error(fmt.Sprintf("SortOrder: value must be one of %v", validSortOrder)),
		),
	)
	if err != nil {
		return err
	}

	return nil
}
