package commentservice

import (
	"context"
	"errors"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"log/slog"
	"time"
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

//go:generate mockery --name Provider
type Provider interface {
	GetComment(ctx context.Context, commentID string) (domain.Comment, error)
	GetPostComments(ctx context.Context, postID string) ([]domain.Comment, error)
}

//go:generate mockery --name Creator
type Creator interface {
	CreateComment(ctx context.Context, comment domain.Comment) (domain.Comment, error)
}

//go:generate mockery --name Updater
type Updater interface {
	UpdateComment(ctx context.Context, comment domain.Comment) (domain.Comment, error)
}

//go:generate mockery --name Deleter
type Deleter interface {
	DeleteComment(ctx context.Context, commentID string) error
}

//go:generate mockery --name UserProvider
type UserProvider interface {
	GetUser(ctx context.Context, id int64) (domain.User, error)
}

func New(config Config) Service {
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

func (s Service) Update(ctx context.Context, dto UpdateCommentDTO) (domain.Comment, error) {
	const op = "service.comment.update"
	log := s.log.With(slog.String("op", op))

	comment, err := s.provider.GetComment(ctx, dto.CommentID)
	if err != nil {
		return domain.Comment{}, handleErr(log, op, err)
	}

	comment.Body = dto.Body
	comment.UpdatedAt = time.Now()

	updatedComment, err := s.updater.UpdateComment(ctx, comment)
	if err != nil {
		return domain.Comment{}, handleErr(log, op, err)
	}

	return updatedComment, nil
}

func (s Service) Delete(ctx context.Context, commentID string) error {
	const op = "service.comment.delete"
	log := s.log.With(slog.String("op", op))

	_, err := s.provider.GetComment(ctx, commentID)
	if err != nil {
		return handleErr(log, op, err)
	}

	err = s.deleter.DeleteComment(ctx, commentID)
	if err != nil {
		return handleErr(log, op, err)
	}

	return nil
}

func (s Service) GetByID(ctx context.Context, commentID string) (domain.Comment, error) {
	const op = "service.comment.get_by_id"
	log := s.log.With(slog.String("op", op))

	comment, err := s.provider.GetComment(ctx, commentID)
	if err != nil {
		return domain.Comment{}, handleErr(log, op, err)
	}

	return comment, nil
}

func (s Service) GetByPostID(ctx context.Context, postID string) ([]domain.Comment, error) {
	const op = "service.comment.get_by_post_id"
	log := s.log.With(slog.String("op", op))

	comments, err := s.provider.GetPostComments(ctx, postID)
	if err != nil {
		return nil, handleErr(log, op, err)
	}

	return comments, nil
}

func handleErr(log *slog.Logger, op string, err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidID):
		return err
	case errors.Is(err, domain.ErrUserNotFound), errors.Is(err, domain.ErrCommentNotFound):
		return err
	case errors.Is(err, domain.ErrInvalidArg):
		return err
	default:
		log.Error(op, logger.Err(err))
		return domain.ErrInternal
	}
}
