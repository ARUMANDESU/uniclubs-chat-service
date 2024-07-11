package domain

import "time"

type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	User      User      `json:"user"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
