package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type contextKey string

const userIDContextKey contextKey = "user_id"

// ContextWithUserID returns a copy of parent context with the authenticated user ID.
func ContextWithUserID(parent context.Context, userID string) context.Context {
	return context.WithValue(parent, userIDContextKey, userID)
}

// UserIDFromContext retrieves the authenticated user ID from context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(userIDContextKey).(string)
	return val, ok && val != ""
}

// RequireAuth middleware verifies the Bearer JWT token and populates the request context with user ID.
func RequireAuth(tm *TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeUnauthorized(w, "Header Authorization diperlukan.")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
				writeUnauthorized(w, "Format header Authorization tidak valid. Gunakan format 'Bearer <token>'.")
				return
			}

			tokenStr := strings.TrimSpace(parts[1])
			claims, err := tm.Verify(tokenStr)
			if err != nil {
				writeUnauthorized(w, "Token tidak valid atau telah kedaluwarsa.")
				return
			}

			ctx := ContextWithUserID(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"code":    "UNAUTHORIZED",
			"message": message,
		},
	})
}
