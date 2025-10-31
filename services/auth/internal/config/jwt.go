package config

import (
	"os"
	"strconv"
	"time"
)

const (
	refreshTokenSecretKeyName = "RT_SECRET"
	accessTokenSecretKeyName  = "AT_SECRET"
	refreshTokenExpansionName = "RT_EXP"
	accessTokenExpansionName  = "AT_EXP"
	authPrefixName            = "AUTH_PREFIX"
)

type JwtConfig interface {
	RefreshTokenSecretKey() []byte
	AccessTokenSecretKey() []byte
	RefreshTokenExpansion() time.Duration
	AccessTokenExpansion() time.Duration
	AuthPrefix() string
}

type jwtConfig struct {
	refreshTokenSecretKey []byte
	accessTokenSecretKey  []byte
	refreshTokenExpansion time.Duration
	accessTokenExpansion  time.Duration
	authPrefix            string
}

func NewJwtConfig() (JwtConfig, error) {
	refreshExpStr := os.Getenv(refreshTokenExpansionName)
	accessExpStr := os.Getenv(accessTokenExpansionName)

	refreshExp, err := strconv.Atoi(refreshExpStr)
	if err != nil {
		refreshExp = 0
	}

	accessExp, err := strconv.Atoi(accessExpStr)
	if err != nil {
		accessExp = 0
	}

	return &jwtConfig{
		refreshTokenSecretKey: []byte(os.Getenv(refreshTokenSecretKeyName)),
		accessTokenSecretKey:  []byte(os.Getenv(accessTokenSecretKeyName)),
		refreshTokenExpansion: time.Duration(refreshExp) * time.Minute,
		accessTokenExpansion:  time.Duration(accessExp) * time.Minute,
		authPrefix:            os.Getenv(authPrefixName),
	}, nil
}

func (c *jwtConfig) RefreshTokenSecretKey() []byte {
	return c.refreshTokenSecretKey
}

func (c *jwtConfig) AccessTokenSecretKey() []byte {
	return c.accessTokenSecretKey
}

func (c *jwtConfig) RefreshTokenExpansion() time.Duration {
	return c.refreshTokenExpansion
}

func (c *jwtConfig) AccessTokenExpansion() time.Duration {
	return c.accessTokenExpansion
}

func (c *jwtConfig) AuthPrefix() string {
	return c.authPrefix
}
