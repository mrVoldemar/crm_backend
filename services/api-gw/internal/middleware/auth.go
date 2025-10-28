package middleware

import (
	"context"
	"net/http"

	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/auth"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/errors"
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
			errors.WriteError(w, http.StatusUnauthorized, "Missing authorization token")
			return
		}

		// Call auth service to validate token
		valid, userID, err := am.authServiceClient.ValidateToken(r.Context(), token)
		if err != nil {
			errors.WriteError(w, http.StatusInternalServerError, "Error validating token")
			return
		}

		if !valid {
			errors.WriteError(w, http.StatusUnauthorized, "Invalid authorization token")
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
