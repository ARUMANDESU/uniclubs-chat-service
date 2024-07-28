package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/storage/mongodb/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Storage) GetUserByID(ctx context.Context, id int64) (domain.User, error) {
	const op = "storage.mongodb.get_user_by_id"

	filter := bson.M{"_id": id}

	var user dao.User
	err := s.commentCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user.ToDomain(), nil
}

func (s *Storage) UpdateUser(ctx context.Context, user domain.User) error {
	const op = "storage.mongodb.update_user"

	filter := bson.M{"user._id": user.ID}

	update := bson.M{
		"$set": bson.M{
			"user": dao.UserFromDomain(user),
		},
	}

	result, err := s.commentCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
