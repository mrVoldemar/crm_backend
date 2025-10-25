package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AuthServiceClient handles communication with the authentication service
type AuthServiceClient struct {
	baseURL string
	client  *http.Client
}

// NewAuthServiceClient creates a new AuthServiceClient
func NewAuthServiceClient(baseURL string) *AuthServiceClient {
	return &AuthServiceClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// ValidateToken validates a token with the auth service
func (asc *AuthServiceClient) ValidateToken(ctx context.Context, token string) (bool, string, error) {
	// In a real implementation, this would make an HTTP request to the auth service
	// to validate the token

	// Example:
	req, err := http.NewRequestWithContext(ctx, "GET", asc.baseURL+"/validate", nil)
	if err != nil {
		return false, "", err
	}

	req.Header.Set("Authorization", token)

	resp, err := asc.client.Do(req)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "", nil
	}

	// Parse response to get user ID
	var result struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, "", fmt.Errorf("error decoding auth service response: %w", err)
	}

	return true, result.UserID, nil

	// Placeholder implementation - remove this when actual service is ready
	// return true, "user123", nil
}
