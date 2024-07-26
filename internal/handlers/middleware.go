package handlers

// func authMiddleware(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
// 			return
// 		}

// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

// 		userID, err := jwt.GetUserID(tokenString, os.Getenv("JWT_SECRET"))
// 		if err != nil {
// 			switch {
// 			case errors.Is(err, domain.ErrTokenIsNotValid),
// 				errors.Is(err, domain.ErrInvalidTokenClaims),
// 				errors.Is(err, domain.ErrUserIDClaimNotFound),
// 				errors.Is(err, domain.ErrTokenIsExpired),
// 				errors.Is(err, domain.ErrUnexpectedSigningMethod):

// 				w.WriteHeader(http.StatusUnauthorized)
// 				w.Write([]byte(err.Error()))
// 			default:
// 				// todo: change log to slog
// 				log.Print("failed to get user id from jwt token", logger.Err(err))
// 				w.WriteHeader(http.StatusInternalServerError)
// 			}
// 			return
// 		}

// 		ctx := r.Context()
// 		newCtx := centrifuge.SetCredentials(ctx, &centrifuge.Credentials{
// 			UserID:   strconv.FormatInt(userID, 10),
// 			ExpireAt: time.Now().Unix() + 60,
// 		})
// 		r = r.WithContext(newCtx)
// 		h.ServeHTTP(w, r)
// 	})
// }
