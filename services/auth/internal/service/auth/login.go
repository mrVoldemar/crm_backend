package auth

import (
	"context"
	"github.com/mrVoldemar/crm_backend/services/auth/internal/model"
	"github.com/mrVoldemar/crm_backend/services/auth/internal/utils"
)

func (s *serv) Login(ctx context.Context, user *model.UserLogin) (string, error) {
	// Лезем в базу или кэш за данными пользователя
	// Сверяем хэши пароля
	//todo:repo-layer дернуть пользователя и пароль

	refreshToken, err := utils.GenerateToken(
		model.UserInfo{
			Username: user.Username,
			// Это пример, в реальности роль должна браться из базы или кэша
			Role: "admin",
		},
		s.jwtConfig.RefreshTokenSecretKey(),
		s.jwtConfig.RefreshTokenExpansion(),
	)

	return refreshToken, err
}
