package auth

import (
	"github.com/mrVoldemar/crm_backend/services/auth/internal/config"
	def "github.com/mrVoldemar/crm_backend/services/auth/internal/service"
)

var _ def.AuthService = (*serv)(nil)

func NewService(jwtConfig config.JwtConfig) *serv {
	return &serv{
		jwtConfig: jwtConfig,
	}
}

type serv struct {
	jwtConfig config.JwtConfig
}
