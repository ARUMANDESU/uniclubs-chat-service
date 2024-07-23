package commentgrpc

import (
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	commentv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/comment"
)

func CommentToProto(c domain.Comment) *commentv1.CommentObject {
	return &commentv1.CommentObject{
		Id:        c.ID,
		PostId:    c.PostID,
		Body:      c.Body,
		CreatedAt: c.CreatedAt.String(),
		UpdatedAt: c.UpdatedAt.String(),
		User:      UserToProto(c.User),
	}
}

func UserToProto(u domain.User) *commentv1.UserObject {
	return &commentv1.UserObject{
		Id:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		AvatarUrl: u.AvatarURL,
	}
}

func PaginationMetadataToProto(p domain.PaginationMetadata) *commentv1.PaginationMetadata {
	return &commentv1.PaginationMetadata{
		CurrentPage:  p.CurrentPage,
		PageSize:     p.PageSize,
		FirstPage:    p.FirstPage,
		LastPage:     p.LastPage,
		TotalRecords: p.TotalRecords,
	}
}

func CommentsToProto(cs []domain.Comment) []*commentv1.CommentObject {
	comments := make([]*commentv1.CommentObject, 0, len(cs))
	for _, c := range cs {
		comments = append(comments, CommentToProto(c))
	}
	return comments
}

func ProtoToFilter(p *commentv1.ListPostCommentsRequest) (domain.Filter, error) {

	var sortBy domain.SortBy
	switch p.SortBy {
	case commentv1.SortBy_SORT_BY_UNSPECIFIED:
		sortBy = domain.SortByUnspecified
	case commentv1.SortBy_SORT_BY_CREATED_AT:
		sortBy = domain.SortByCreatedAt
	case commentv1.SortBy_SORT_BY_UPDATED_AT:
		sortBy = domain.SortByUpdatedAt
	default:
		sortBy = domain.SortByUnspecified
	}

	var sortOrder domain.SortOrder
	switch p.SortOrder {
	case commentv1.SortOrder_SORT_ORDER_ASC:
		sortOrder = domain.SortOrderAsc
	case commentv1.SortOrder_SORT_ORDER_DESC:
		sortOrder = domain.SortOrderDesc
	default:
		sortOrder = domain.SortOrderDesc
	}

	filter, err := domain.NewFilter(
		domain.WithPage(p.Page),
		domain.WithPageSize(p.PageSize),
		domain.WithSortBy(sortBy),
		domain.WithSortOrder(sortOrder),
	)
	if err != nil {
		return domain.Filter{}, err
	}

	return *filter, nil
}
