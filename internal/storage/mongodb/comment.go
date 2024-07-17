package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/storage"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/storage/mongodb/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Storage) GetComment(ctx context.Context, id string) (domain.Comment, error) {
	const op = "storage.mongodb.get_comment"

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return domain.Comment{}, fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return domain.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	var comment dao.Comment
	err = s.commentCollection.FindOne(ctx, objectID).Decode(&comment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Comment{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return domain.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment.ToDomain(), nil
}

func (s *Storage) GetPostComments(ctx context.Context, postID string) ([]domain.Comment, error) {
	const op = "storage.mongodb.get_post_comments"

	cursor, err := s.commentCollection.Find(ctx, bson.M{"": postID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var comments []dao.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.CommentsToDomain(comments), nil
}

func (s *Storage) CreateComment(ctx context.Context, domainComment domain.Comment) (domain.Comment, error) {
	const op = "storage.mongodb.create_comment"

	comment := dao.CommentFromDomain(domainComment)
	comment.ID = primitive.NewObjectID()

	result, err := s.commentCollection.InsertOne(ctx, comment)
	if err != nil {
		return domain.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.commentCollection.FindOne(ctx, result.InsertedID).Decode(&comment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Comment{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return domain.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment.ToDomain(), nil
}

func (s *Storage) DeleteComment(ctx context.Context, id string) error {
	const op = "storage.mongodb.delete_comment"

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.commentCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateComment(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
	//TODO implement me
	panic("implement me")
}
