package service

import (
	"context"

	"github.com/mrVoldemar/crm_backend/services/auth/internal/model"
)

type AuthService interface {
	Login(ctx context.Context, user *model.UserLogin) (string, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type AccessService interface {
	Check(ctx context.Context, endpoint string) error
}
