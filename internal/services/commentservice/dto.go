package commentservice

type CreateCommentDTO struct {
	PostID string `json:"post_id"`
	Body   string `json:"body"`
	UserID int64  `json:"user_id"`
}

type UpdateCommentDTO struct {
	CommentID string `json:"comment_id"`
	Body      string `json:"body"`
}
