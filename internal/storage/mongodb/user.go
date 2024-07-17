package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/storage"
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
			return domain.User{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user.ToDomain(), nil
}

func (s *Storage) SaveUser(ctx context.Context, domainUser domain.User) error {
	const op = "storage.mongodb.save_user"

	user := dao.UserFromDomain(domainUser)

	_, err := s.commentCollection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
