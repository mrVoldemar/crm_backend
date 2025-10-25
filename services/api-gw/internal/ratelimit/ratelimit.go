package ratelimit

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a simple rate limiter
type RateLimiter struct {
	limit    int64
	interval time.Duration
	tokens   map[string]*clientTokens
	mutex    sync.Mutex
}

type clientTokens struct {
	count      int64
	lastRefill time.Time
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(limit int64, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:    limit,
		interval: interval,
		tokens:   make(map[string]*clientTokens),
	}
}

// IsAllowed checks if the client is allowed to make a request
func (rl *RateLimiter) IsAllowed(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	client, exists := rl.tokens[clientIP]
	if !exists {
		rl.tokens[clientIP] = &clientTokens{
			count:      rl.limit - 1,
			lastRefill: now,
		}
		return true
	}

	// Refill tokens based on time passed
	if now.Sub(client.lastRefill) >= rl.interval {
		client.count = rl.limit
		client.lastRefill = now
	}

	// Check if tokens are available
	if client.count > 0 {
		client.count--
		return true
	}

	return false
}

// RateLimitMiddleware applies rate limiting to HTTP requests
func (rl *RateLimiter) RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr

		if !rl.IsAllowed(clientIP) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	}
}
