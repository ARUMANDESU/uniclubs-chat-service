package commentservice

import (
	"context"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"log/slog"
	"time"
)

type Config struct {
	Logger   *slog.Logger
	Provider Provider
	Creator  Creator
	Updater  Updater
	Deleter  Deleter
}

type Service struct {
	log      *slog.Logger
	provider Provider
	creator  Creator
	updater  Updater
	deleter  Deleter
}

type Provider interface {
	GetComment(ctx context.Context, commentID string) (domain.Comment, error)
	GetPostComments(ctx context.Context, postID string) ([]domain.Comment, error)
}

type Creator interface {
	CreateComment(ctx context.Context, comment CreateCommentDTO) (domain.Comment, error)
}

type Updater interface {
	UpdateComment(ctx context.Context, comment domain.Comment) (domain.Comment, error)
}

type Deleter interface {
	DeleteComment(ctx context.Context, commentID string) error
}

func NewComment(config Config) Service {
	return Service{
		log:      config.Logger,
		provider: config.Provider,
		creator:  config.Creator,
		updater:  config.Updater,
		deleter:  config.Deleter,
	}
}

func (c Service) Create(ctx context.Context, comment CreateCommentDTO) (domain.Comment, error) {
	//TODO implement me

	return domain.Comment{
		ID:        "1",
		PostID:    "1",
		Body:      "Hello World",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (c Service) Update(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (c Service) Delete(ctx context.Context, commentID string) error {
	//TODO implement me
	panic("implement me")
}
