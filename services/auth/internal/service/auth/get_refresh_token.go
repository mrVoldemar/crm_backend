package auth

import (
	"context"
	"github.com/mrVoldemar/crm_backend/services/auth/internal/model"
	"github.com/mrVoldemar/crm_backend/services/auth/internal/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serv) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, s.jwtConfig.RefreshTokenSecretKey())
	if err != nil {
		return "", status.Errorf(codes.Aborted, "invalid refresh token")
	}

	// Можем слазать в базу или в кэш за доп данными пользователя

	newRefreshToken, err := utils.GenerateToken(model.UserInfo{
		Username: claims.Username,
		// Это пример, в реальности роль должна браться из базы или кэша
		Role: "admin",
	},
		s.jwtConfig.RefreshTokenSecretKey(),
		s.jwtConfig.RefreshTokenExpansion(),
	)

	return newRefreshToken, err
}
