package mongodb

import (
	"context"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/config"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/services/commentservice"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Storage struct {
	client            *mongo.Client
	commentCollection *mongo.Collection
}

// NewStorage creates a new MongoDB storage instance
func NewStorage(ctx context.Context, cfg config.MongoDB) (Storage, error) {
	const op = "storage.mongodb.new"

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return Storage{}, fmt.Errorf("%s: %w", op, err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, cfg.PingTimeout)
	defer cancel()
	if err = client.Ping(pingCtx, readpref.Primary()); err != nil {
		return Storage{}, fmt.Errorf("%s: %w", op, err)
	}

	db := client.Database(cfg.DatabaseName)
	commentsCollection := db.Collection("comments")

	return Storage{client: client, commentCollection: commentsCollection}, nil
}

func (s *Storage) Stop(ctx context.Context) error {
	const op = "storage.mongodb.close"

	if err := s.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetComment(ctx context.Context, commentID string) (domain.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetPostComments(ctx context.Context, postID string) ([]domain.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) CreateComment(ctx context.Context, comment commentservice.CreateCommentDTO) (domain.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) DeleteComment(ctx context.Context, commentID string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) UpdateComment(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
	//TODO implement me
	panic("implement me")
}
