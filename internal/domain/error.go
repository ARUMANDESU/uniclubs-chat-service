package domain

import "errors"

var (
	ErrInvalidID  = errors.New("invalid id")
	ErrInvalidArg = errors.New("invalid argument")
	ErrInternal   = errors.New("internal error")

	ErrUserNotFound = errors.New("user not found")
)

var (
	ErrCommentNotFound = errors.New("comment not found")
)
