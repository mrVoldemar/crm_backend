package handlers

import (
	"net/http"

	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/middleware"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/proxy"
)

// APIHandler handles API requests
type APIHandler struct {
	authMiddleware *middleware.AuthMiddleware
	employeeProxy  *proxy.ServiceProxy
	clientProxy    *proxy.ServiceProxy
	// Add other service proxies as needed
}

// NewAPIHandler creates a new APIHandler
func NewAPIHandler(authMiddleware *middleware.AuthMiddleware) (*APIHandler, error) {
	// Initialize service proxies
	employeeProxy, err := proxy.NewServiceProxy("http://employee-service:8080")
	if err != nil {
		return nil, err
	}

	clientProxy, err := proxy.NewServiceProxy("http://client-service:8080")
	if err != nil {
		return nil, err
	}

	return &APIHandler{
		authMiddleware: authMiddleware,
		employeeProxy:  employeeProxy,
		clientProxy:    clientProxy,
	}, nil
}

// RegisterRoutes registers all API routes
func (h *APIHandler) RegisterRoutes(mux *http.ServeMux) {
	// Employee service routes
	mux.HandleFunc("/api/employees", h.authMiddleware.AuthRequired(
		h.employeeProxy.WithAuth(h.employeeProxy.ProxyRequest)))
	mux.HandleFunc("/api/employees/", h.authMiddleware.AuthRequired(
		h.employeeProxy.WithAuth(h.employeeProxy.ProxyRequest)))

	// Client service routes
	mux.HandleFunc("/api/clients", h.authMiddleware.AuthRequired(
		h.clientProxy.WithAuth(h.clientProxy.ProxyRequest)))
	mux.HandleFunc("/api/clients/", h.authMiddleware.AuthRequired(
		h.clientProxy.WithAuth(h.clientProxy.ProxyRequest)))

	// Health check endpoint
	mux.HandleFunc("/health", h.healthCheck)
}

func (h *APIHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API Gateway is running"))
}
