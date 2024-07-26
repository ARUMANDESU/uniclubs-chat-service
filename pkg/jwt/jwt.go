package jwt

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

// GetUserID extracts the user_id from the provided JWT token
// returns the user_id, expiresAt and an error if the token is not valid
func GetUserID(tokenString, secret string) (int64, int64, error) {
	const op = "jwt.GetUserID"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret used to sign the token
		return []byte(secret), nil
	})
	if err != nil {
		errorMessage := err.Error()
		switch {
		case strings.Contains(errorMessage, "token is expired"):
			return 0, 0, domain.ErrTokenIsExpired
		case strings.Contains(errorMessage, "token contains an invalid number of segments"):
			return 0, 0, domain.ErrTokenIsNotValid
		case strings.Contains(errorMessage, "unexpected signing method"):
			return 0, 0, domain.ErrUnexpectedSigningMethod
		}
		return 0, 0, fmt.Errorf("%s: %w", op, err)
	}

	// Check if the token is valid
	if !token.Valid {
		return 0, 0, domain.ErrTokenIsNotValid
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, 0, domain.ErrInvalidTokenClaims
	}

	// Extract user_id from claims
	expTime, ok := claims["exp"].(float64)
	slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})).Debug(fmt.Sprintf("%v", expTime))
	if !ok {
		return 0, 0, domain.ErrInvalidTokenClaims
	}

	if time.Now().Unix() > int64(expTime) {
		return 0, 0, domain.ErrTokenIsExpired
	}

	// Extract user_id from claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, 0, domain.ErrUserIDClaimNotFound
	}

	return int64(userID), int64(expTime), nil
}
