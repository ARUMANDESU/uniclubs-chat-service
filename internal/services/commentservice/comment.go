package commentservice

import (
	"context"
	"errors"
	userclient "github.com/ARUMANDESU/uniclubs-comments-service/internal/client/user"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/storage"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"log/slog"
	"time"
)

var (
	ErrInvalidID  = errors.New("invalid id")
	ErrNotFound   = errors.New("not found")
	ErrInvalidArg = errors.New("invalid argument")
)

type Config struct {
	Logger       *slog.Logger
	Provider     Provider
	Creator      Creator
	Updater      Updater
	Deleter      Deleter
	UserProvider UserProvider
}

type Service struct {
	log          *slog.Logger
	provider     Provider
	creator      Creator
	updater      Updater
	deleter      Deleter
	userProvider UserProvider
}

type Provider interface {
	GetComment(ctx context.Context, commentID string) (domain.Comment, error)
	GetPostComments(ctx context.Context, postID string) ([]domain.Comment, error)
}

type Creator interface {
	CreateComment(ctx context.Context, comment domain.Comment) (domain.Comment, error)
}

type Updater interface {
	UpdateComment(ctx context.Context, comment domain.Comment) (domain.Comment, error)
}

type Deleter interface {
	DeleteComment(ctx context.Context, commentID string) error
}

type UserProvider interface {
	GetUser(ctx context.Context, id int64) (domain.User, error)
}

func NewComment(config Config) Service {
	return Service{
		log:          config.Logger,
		provider:     config.Provider,
		creator:      config.Creator,
		updater:      config.Updater,
		deleter:      config.Deleter,
		userProvider: config.UserProvider,
	}
}

func (s Service) Create(ctx context.Context, comment CreateCommentDTO) (domain.Comment, error) {
	const op = "service.comment.create"
	log := s.log.With(slog.String("op", op))

	user, err := s.userProvider.GetUser(ctx, comment.UserID)
	if err != nil {
		return domain.Comment{}, handleErr(log, op, err)
	}

	createdComment, err := s.creator.CreateComment(ctx, domain.Comment{
		PostID:    comment.PostID,
		User:      user,
		Body:      comment.Body,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return domain.Comment{}, handleErr(log, op, err)
	}

	return createdComment, nil
}

func (s Service) Update(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) Delete(ctx context.Context, commentID string) error {
	//TODO implement me
	panic("implement me")
}

func handleErr(log *slog.Logger, op string, err error) error {
	switch {
	case errors.Is(err, storage.ErrInvalidID):
		return ErrInvalidID
	case errors.Is(err, storage.ErrNotFound), errors.Is(err, userclient.ErrUserNotFound):
		return ErrNotFound
	case errors.Is(err, userclient.ErrInvalidArg):
		return ErrInvalidArg
	default:
		log.Error(op, logger.Err(err))
		return err
	}
}
