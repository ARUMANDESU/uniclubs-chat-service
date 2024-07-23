package commentgrpc

import (
	"context"
	"errors"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	commentv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/comment"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	commentv1.UnimplementedCommentServer
	commentService CommentService
}

//go:generate mockery --name CommentService --output ./mocks
type CommentService interface {
	GetByID(ctx context.Context, id string) (domain.Comment, error)
	ListByPostID(ctx context.Context, postID string, filter domain.Filter) ([]domain.Comment, domain.PaginationMetadata, error)
}

func Register(gRPC *grpc.Server, server commentv1.CommentServer) {
	commentv1.RegisterCommentServer(gRPC, server)
}

func NewServer(service CommentService) Server {
	return Server{
		commentService: service,
	}
}

func (s *Server) GetCommentByID(ctx context.Context, request *commentv1.GetCommentByIDRequest) (*commentv1.GetCommentByIDResponse, error) {
	err := validation.ValidateStruct(request,
		validation.Field(&request.Id, validation.Required.Error("comment id is required")),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	comment, err := s.commentService.GetByID(ctx, request.Id)
	if err != nil {
		return nil, handleErr(err)
	}

	return &commentv1.GetCommentByIDResponse{
		Comment: CommentToProto(comment),
	}, nil
}

func (s *Server) ListPostComments(ctx context.Context, request *commentv1.ListPostCommentsRequest) (*commentv1.ListPostCommentsResponse, error) {
	err := validateListPostComments(request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	filter, err := ProtoToFilter(request)
	if err != nil {
		return nil, handleErr(err)
	}

	comments, pagination, err := s.commentService.ListByPostID(ctx, request.PostId, filter)
	if err != nil {
		return nil, handleErr(err)
	}

	return &commentv1.ListPostCommentsResponse{
		Comments: CommentsToProto(comments),
		Metadata: PaginationMetadataToProto(pagination),
	}, nil
}

func handleErr(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidID), errors.Is(err, domain.ErrInvalidArg):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrCommentNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
