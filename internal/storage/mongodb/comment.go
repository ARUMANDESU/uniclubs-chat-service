package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/storage/mongodb/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Storage) GetComment(ctx context.Context, id string) (domain.Comment, error) {
	const op = "storage.mongodb.get_comment"

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return domain.Comment{}, domain.ErrInvalidID
		}
		return domain.Comment{}, fmt.Errorf("%s: failed to convert id to ObjectID: %w", op, err)
	}

	var comment dao.Comment
	err = s.commentCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&comment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Comment{}, domain.ErrCommentNotFound
		}
		return domain.Comment{}, fmt.Errorf("%s: failed to find document: %w", op, err)
	}

	return comment.ToDomain(), nil
}

func (s *Storage) ListPostComments(ctx context.Context, postID string, filters domain.Filter) (
	[]domain.Comment,
	domain.PaginationMetadata,
	error,
) {
	const op = "storage.mongodb.list_post_comments"

	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, domain.PaginationMetadata{}, domain.ErrInvalidID
		}
		return nil, domain.PaginationMetadata{}, fmt.Errorf("%s: failed to convert postID to ObjectID: %w", op, err)
	}

	query := bson.M{"post_id": objectID}

	totalRecords, err := s.commentCollection.CountDocuments(ctx, query)
	if err != nil {
		return nil, domain.PaginationMetadata{}, fmt.Errorf("%s: failed to count documents: %w", op, err)
	}
	if totalRecords == 0 {
		return nil, domain.PaginationMetadata{}, domain.ErrCommentNotFound
	}

	if filters.SortBy == "" {
		filters.SortBy = "created_at"
	}
	sort := bson.M{string(filters.SortBy): filters.SortOrder.Mongo()}

	opts := options.Find()
	opts.SetSort(sort)
	opts.SetSkip(int64(filters.Offset()))
	opts.SetLimit(int64(filters.Limit()))

	cursor, err := s.commentCollection.Find(ctx, query, opts)
	if err != nil {
		return nil, domain.PaginationMetadata{}, fmt.Errorf("%s: failed to find documents: %w", op, err)
	}

	var comments []dao.Comment
	err = cursor.All(ctx, &comments)
	if err != nil {
		return nil, domain.PaginationMetadata{}, fmt.Errorf("%s: failed to decode documents: %w", op, err)
	}

	paginationMetadata := domain.CalculatePaginationMetadata(int32(totalRecords), filters.Page, filters.PageSize)

	return dao.CommentsToDomain(comments), paginationMetadata, nil
}

func (s *Storage) CreateComment(ctx context.Context, domainComment domain.Comment) (domain.Comment, error) {
	const op = "storage.mongodb.create_comment"

	comment, err := dao.CommentFromDomain(domainComment)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return domain.Comment{}, domain.ErrInvalidID
		}
		return domain.Comment{}, fmt.Errorf("%s: failed to convert domain comment to dao: %w", op, err)
	}

	result, err := s.commentCollection.InsertOne(ctx, comment)
	if err != nil {
		return domain.Comment{}, fmt.Errorf("%s: failed to insert document: %w", op, err)
	}

	err = s.commentCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&comment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Comment{}, domain.ErrCommentNotFound
		}
		return domain.Comment{}, fmt.Errorf("%s: failed to find document: %w", op, err)
	}

	return comment.ToDomain(), nil
}

func (s *Storage) DeleteComment(ctx context.Context, id string) error {
	const op = "storage.mongodb.delete_comment"

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return domain.ErrInvalidID
		}
		return fmt.Errorf("%s: failed to convert id to ObjectID: %w", op, err)
	}

	_, err = s.commentCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("%s: failed to delete document: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateComment(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
	const op = "storage.mongodb.update_comment"

	objectID, err := primitive.ObjectIDFromHex(comment.ID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return domain.Comment{}, domain.ErrInvalidID
		}
		return domain.Comment{}, fmt.Errorf("%s failed to convert id to ObjectID: %w", op, err)
	}

	commentToUpdate, _ := dao.CommentFromDomain(comment)
	commentToUpdate.ID = objectID

	_, err = s.commentCollection.ReplaceOne(ctx, bson.M{"_id": objectID}, commentToUpdate)
	if err != nil {
		return domain.Comment{}, fmt.Errorf("%s failed to update document: %w", op, err)
	}

	return comment, nil
}
