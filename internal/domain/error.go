package domain

import "errors"

var (
	ErrInvalidID    = errors.New("invalid id")
	ErrInvalidArg   = errors.New("invalid argument")
	ErrInternal     = errors.New("internal error")
	ErrUnauthorized = errors.New("unauthorized")

	ErrTokenIsNotValid         = errors.New("token is not valid")
	ErrInvalidTokenClaims      = errors.New("invalid token claims")
	ErrUserIDClaimNotFound     = errors.New("user_id claim not found or invalid")
	ErrTokenIsExpired          = errors.New("token is expired")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)

var (
	ErrUserNotFound = errors.New("user not found")
)

var (
	ErrCommentNotFound = errors.New("comment not found")
)
