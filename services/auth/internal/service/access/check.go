package access

import (
	"context"
	"github.com/mrVoldemar/crm_backend/services/auth/internal/utils"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

func (s *serv) Check(ctx context.Context, endpoint string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], s.jwtConfig.AuthPrefix()) {
		return errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], s.jwtConfig.AuthPrefix())

	claims, err := utils.VerifyToken(accessToken, s.jwtConfig.AccessTokenSecretKey())
	if err != nil {
		return errors.New("access token is invalid")
	}

	accessibleMap, err := s.accessibleRoles(ctx)
	if err != nil {
		return errors.New("failed to get accessible roles")
	}

	role, ok := accessibleMap[endpoint]
	if !ok {
		return nil
	}

	if role == claims.Role {
		return nil
	}

	return errors.New("access denied")
}
