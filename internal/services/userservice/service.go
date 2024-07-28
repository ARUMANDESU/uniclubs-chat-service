package userservice

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
)

type Service struct {
	log *slog.Logger
	// dbProvider is the primary provider
	dbProvider UserProvider
	// grpcProvider is remote user microservice, it is the secondary provider called when the primary provider fails
	grpcProvider UserGRPCProvider
}

//go:generate mockery --name UserProvider
type UserProvider interface {
	GetUserByID(ctx context.Context, id int64) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) error
}

//go:generate mockery --name UserGRPCProvider
type UserGRPCProvider interface {
	GetUserByID(ctx context.Context, id int64) (domain.User, error)
}

func New(log *slog.Logger, primaryProvider UserProvider, secondaryProvider UserGRPCProvider) Service {
	return Service{
		log:          log,
		dbProvider:   primaryProvider,
		grpcProvider: secondaryProvider,
	}
}

func (s *Service) GetUser(ctx context.Context, id int64) (domain.User, error) {
	const op = "service.user.get_user"
	log := s.log.With(slog.String("op", op))

	user, err := s.dbProvider.GetUserByID(ctx, id)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			log.Warn("db provider failed", logger.Err(err))
		}

		user, err = s.grpcProvider.GetUserByID(ctx, id)
		if err != nil {
			return domain.User{}, handleErr(log, op, err)
		}
	}

	return user, nil
}

func (s *Service) Update(ctx context.Context, user domain.User) error {
	const op = "service.user.update"
	log := s.log.With(slog.String("op", op))

	err := s.dbProvider.UpdateUser(ctx, user)
	if err != nil {
		return handleErr(log, op, err)
	}

	return nil
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
