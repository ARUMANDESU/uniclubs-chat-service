package userservice

import (
	"context"
	"errors"
	userclient "github.com/ARUMANDESU/uniclubs-comments-service/internal/client/user"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"log/slog"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidArg   = errors.New("invalid argument")
)

type Service struct {
	log               *slog.Logger
	primaryProvider   UserProvider
	secondaryProvider UserProvider
	saver             UserSaver
}

//go:generate mockery --name UserProvider
type UserProvider interface {
	GetUserByID(ctx context.Context, id int64) (domain.User, error)
}

//go:generate mockery --name UserSaver
type UserSaver interface {
	SaveUser(ctx context.Context, user domain.User) error
}

func New(log *slog.Logger, primaryProvider, secondaryProvider UserProvider, saver UserSaver) Service {
	return Service{
		log:               log,
		primaryProvider:   primaryProvider,
		secondaryProvider: secondaryProvider,
		saver:             saver,
	}
}

func (s *Service) GetUser(ctx context.Context, id int64) (domain.User, error) {
	const op = "service.user.get_user"
	log := s.log.With(slog.String("op", op))

	user, err := s.primaryProvider.GetUserByID(ctx, id)
	if err != nil {
		if !errors.Is(err, userclient.ErrUserNotFound) {
			log.Error("primary provider failed", logger.Err(err))
		}

		user, err = s.secondaryProvider.GetUserByID(ctx, id)
		if err != nil {
			return domain.User{}, handleErr(log, op, err)
		}
		go func() {
			saveCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			if err := s.saver.SaveUser(saveCtx, user); err != nil {
				log.Error("failed to save user", logger.Err(err))
			}

			<-saveCtx.Done()
			if saveCtx.Err() != nil {
				log.Error("failed to save user", logger.Err(saveCtx.Err()))
			}
		}()
	}

	return user, nil
}

func handleErr(log *slog.Logger, op string, err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, userclient.ErrUserNotFound):
		return ErrUserNotFound
	case errors.Is(err, userclient.ErrInvalidArg):
		return ErrInvalidArg
	default:
		log.Error(op, logger.Err(err))
		return err
	}
}
