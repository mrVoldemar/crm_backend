package middleware

import (
	"context"
	"net/http"

	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/auth"
)

// AuthMiddleware handles authentication and authorization
type AuthMiddleware struct {
	authServiceClient *auth.AuthServiceClient
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(authServiceURL string) *AuthMiddleware {
	client := auth.NewAuthServiceClient(authServiceURL)
	return &AuthMiddleware{
		authServiceClient: client,
	}
}

// AuthRequired is a middleware function that checks if the user is authenticated
func (am *AuthMiddleware) AuthRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// In a real implementation, we would make a call to the auth service
		// to validate the token provided in the Authorization header
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		// Call auth service to validate token
		valid, userID, err := am.authServiceClient.ValidateToken(r.Context(), token)
		if err != nil {
			http.Error(w, "Error validating token", http.StatusInternalServerError)
			return
		}

		if !valid {
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// validateToken would call the auth service to validate the token
// func (am *AuthMiddleware) validateToken(token string) (bool, string, error) {
// 	// This would make an HTTP/gRPC call to the auth service
// 	// For now, returning a placeholder
// 	return true, "user123", nil
// }
