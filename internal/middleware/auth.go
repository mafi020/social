package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/mafi020/social/internal/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.JSONErrorResponse(w, http.StatusUnauthorized, map[string]string{"message": "Missing Authorization header"})
			return
		}

		// Check format "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.JSONErrorResponse(w, http.StatusUnauthorized, map[string]string{"message": "Invalid Authorization header format"})
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(parts[1])
		if err != nil {

			utils.JSONErrorResponse(w, http.StatusUnauthorized, map[string]string{"message": "Invalid or expired token"})
			return
		}

		// Put user ID into context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper to get userID from context in handlers
func GetAuthUserIDFromContext(r *http.Request) int64 {
	if userID, ok := r.Context().Value(UserIDKey).(int64); ok {
		return userID
	}
	return 0
}
