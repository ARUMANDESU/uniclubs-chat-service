package dao

import (
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"time"
)

type Comment struct {
	ID        string    `json:"id" bson:"_id"`
	PostID    string    `json:"post_id" bson:"post_id"`
	User      User      `json:"user" bson:"user"`
	Body      string    `json:"body" bson:"body"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (c *Comment) ToDomain() domain.Comment {
	if c == nil {
		return domain.Comment{}
	}

	return domain.Comment{
		ID:        c.ID,
		PostID:    c.PostID,
		User:      c.User.ToDomain(),
		Body:      c.Body,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func CommentFromDomain(d domain.Comment) Comment {
	return Comment{
		ID:        d.ID,
		PostID:    d.PostID,
		User:      UserFromDomain(d.User),
		Body:      d.Body,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func CommentsToDomain(comments []Comment) []domain.Comment {
	domainComments := make([]domain.Comment, 0, len(comments))
	for _, comment := range comments {
		domainComments = append(domainComments, comment.ToDomain())
	}
	return domainComments
}
