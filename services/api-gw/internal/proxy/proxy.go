package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ServiceProxy handles proxying requests to downstream services
type ServiceProxy struct {
	targetServiceURL string
	proxy            *httputil.ReverseProxy
}

// NewServiceProxy creates a new ServiceProxy
func NewServiceProxy(targetServiceURL string) (*ServiceProxy, error) {
	parsedURL, err := url.Parse(targetServiceURL)
	if err != nil {
		return nil, fmt.Errorf("invalid target service URL: %w", err)
	}

	return &ServiceProxy{
		targetServiceURL: targetServiceURL,
		proxy:            httputil.NewSingleHostReverseProxy(parsedURL),
	}, nil
}

// ProxyRequest handles proxying the request to the target service
func (sp *ServiceProxy) ProxyRequest(w http.ResponseWriter, r *http.Request) {
	// Parse the target URL to properly configure the reverse proxy
	targetURL, err := url.Parse(sp.targetServiceURL)
	if err != nil {
		http.Error(w, "Invalid target service URL", http.StatusInternalServerError)
		return
	}

	// Update request URL for proxying
	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = targetURL.Host

	// Proxy the request to the target service
	sp.proxy.ServeHTTP(w, r)
}

// WithAuth adds user information to the proxied request
func (sp *ServiceProxy) WithAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add user info from context to headers
		userID := r.Context().Value("userID")
		if userID != nil {
			r.Header.Set("X-User-ID", userID.(string))
		}

		next.ServeHTTP(w, r)
	}
}
