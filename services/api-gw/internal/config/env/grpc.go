package env

import (
	"net"
	"os"

	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/config"
	"github.com/pkg/errors"
)

var _ config.GRPCConfig = (*grpcConfig)(nil)

const (
	grpcHostEnvName       = "GRPC_HOST"
	grpcPortEnvName       = "GRPC_PORT"
	authServiceURLEnvName = "AUTH_SERVICE_URL"
)

type grpcConfig struct {
	host           string
	port           string
	authServiceURL string
}

func NewGRPCConfig() (*grpcConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	authServiceURL := os.Getenv(authServiceURLEnvName)
	if len(authServiceURL) == 0 {
		authServiceURL = "http://auth-service:8080" // default value
	}

	return &grpcConfig{
		host:           host,
		port:           port,
		authServiceURL: authServiceURL,
	}, nil
}

func (cfg *grpcConfig) GRPCAddress() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *grpcConfig) AuthServiceURL() string {
	return cfg.authServiceURL
}
