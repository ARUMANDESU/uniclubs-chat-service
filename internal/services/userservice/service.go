package userservice

import (
	"context"
	"errors"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"log/slog"
)

type Service struct {
	log               *slog.Logger
	primaryProvider   UserProvider
	secondaryProvider UserProvider
}

//go:generate mockery --name UserProvider
type UserProvider interface {
	GetUserByID(ctx context.Context, id int64) (domain.User, error)
}

func New(log *slog.Logger, primaryProvider, secondaryProvider UserProvider) Service {
	return Service{
		log:               log,
		primaryProvider:   primaryProvider,
		secondaryProvider: secondaryProvider,
	}
}

func (s *Service) GetUser(ctx context.Context, id int64) (domain.User, error) {
	const op = "service.user.get_user"
	log := s.log.With(slog.String("op", op))

	user, err := s.primaryProvider.GetUserByID(ctx, id)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			log.Warn("primary provider failed", logger.Err(err))
		}

		user, err = s.secondaryProvider.GetUserByID(ctx, id)
		if err != nil {
			return domain.User{}, handleErr(log, op, err)
		}
	}

	return user, nil
}

func handleErr(log *slog.Logger, op string, err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return err
	case errors.Is(err, domain.ErrInvalidArg):
		return err
	default:
		log.Error(op, logger.Err(err))
		return domain.ErrInternal
	}
}
