package handlers

import (
	"net/http"

	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/config/services"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/errors"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/middleware"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/proxy"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/ratelimit"
)

// APIHandler handles API requests
type APIHandler struct {
	authMiddleware *middleware.AuthMiddleware
	rateLimiter    *ratelimit.RateLimiter
	serviceProxies map[string]*proxy.ServiceProxy
}

// NewAPIHandler creates a new APIHandler
func NewAPIHandler(authMiddleware *middleware.AuthMiddleware, rateLimiter *ratelimit.RateLimiter, serviceConfig *services.Config) (*APIHandler, error) {
	// Initialize service proxies from config
	serviceProxies := make(map[string]*proxy.ServiceProxy)

	for _, svc := range serviceConfig.Services {
		proxy, err := proxy.NewServiceProxy(svc.URL)
		if err != nil {
			return nil, err
		}
		serviceProxies[svc.Name] = proxy
	}

	return &APIHandler{
		authMiddleware: authMiddleware,
		rateLimiter:    rateLimiter,
		serviceProxies: serviceProxies,
	}, nil
}

// RegisterRoutes registers all API routes
func (h *APIHandler) RegisterRoutes(mux *http.ServeMux) {
	// Register routes for each service
	for serviceName, serviceProxy := range h.serviceProxies {
		// Service routes
		mux.HandleFunc("/api/"+serviceName+"s", errors.ErrorHandler(h.rateLimiter.RateLimitMiddleware(h.authMiddleware.AuthRequired(
			serviceProxy.WithAuth(serviceProxy.ProxyRequest)))))
		mux.HandleFunc("/api/"+serviceName+"s/", errors.ErrorHandler(h.rateLimiter.RateLimitMiddleware(h.authMiddleware.AuthRequired(
			serviceProxy.WithAuth(serviceProxy.ProxyRequest)))))
	}

	// Health check endpoint (no rate limiting or auth required)
	mux.HandleFunc("/health", errors.ErrorHandler(h.healthCheck))
}

func (h *APIHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API Gateway is running"))
}
