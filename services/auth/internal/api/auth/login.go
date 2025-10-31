package auth

import (
	"context"
	"github.com/mrVoldemar/crm_backend/services/auth/internal/converter"
	descAuth "github.com/mrVoldemar/crm_backend/services/auth/pkg/auth_v1"
	"time"

	"github.com/pkg/errors"
)

const (
	grpcPort   = 50050
	authPrefix = "Bearer "

	refreshTokenSecretKey = "W4/X+LLjehdxptt4YgGFCvMpq5ewptpZZYRHY6A72g0="
	accessTokenSecretKey  = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="

	refreshTokenExpiration = 60 * time.Minute
	accessTokenExpiration  = 5 * time.Minute
)

func (i *Implementation) Login(ctx context.Context, req *descAuth.LoginRequest) (*descAuth.LoginResponse, error) {

	refreshToken, err := i.authService.Login(ctx, converter.FromRequestToLoginModel(req))

	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &descAuth.LoginResponse{RefreshToken: refreshToken}, nil
}
