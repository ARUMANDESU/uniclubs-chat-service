package mongodb

import (
	"context"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/config"
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
