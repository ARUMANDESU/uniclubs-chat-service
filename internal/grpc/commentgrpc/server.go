package commentgrpc

import (
	"context"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	commentv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/comment"
	"google.golang.org/grpc"
)

type Server struct {
	commentv1.UnimplementedCommentServer
	commentService CommentService
}

type CommentService interface {
	GetByID(ctx context.Context, id string) (domain.Comment, error)
	ListByPostID(ctx context.Context, postID string, filter domain.Filter) ([]domain.Comment, domain.PaginationMetadata, error)
}

func Register(gRPC *grpc.Server, server Server) {
	commentv1.RegisterCommentServer(gRPC, &server)
}

func NewServer(service CommentService) Server {
	return Server{
		commentService: service,
	}
}

func (s *Server) GetCommentByID(ctx context.Context, request *commentv1.GetCommentByIDRequest) (*commentv1.GetCommentByIDResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) ListPostComments(ctx context.Context, request *commentv1.ListPostCommentsRequest) (*commentv1.ListPostCommentsResponse, error) {
	//TODO implement me
	panic("implement me")
}
