package env

import (
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	refreshTokenExpTimeName   = "REFRESH_TOKEN_EXPIRATION_TIME"
	refreshTokenSecretKeyName = "REFRESH_TOKEN_SECRET_KEY"

	accessTokenExpTimeName   = "ACCESS_TOKEN_EXPIRATION_TIME"
	accessTokenSecretKeyName = "ACCESS_TOKEN_SECRET_KEY"
)

type TokenConfig interface {
	ExpTime() time.Duration
	Secret() []byte
}

type tokenConfig struct {
	expTime time.Duration
	secret  []byte
}

func NewRefreshTokenConfig() (TokenConfig, error) {
	expTimeStr := os.Getenv(refreshTokenExpTimeName)
	if expTimeStr == "" {
		return nil, errors.New("refresh token expiration time not found")
	}
	expTime, err := strconv.Atoi(expTimeStr)
	if err != nil {
		return nil, errors.New("unable to get refresh token expiration time")
	}

	secret := os.Getenv(refreshTokenSecretKeyName)
	if secret == "" {
		return nil, errors.New("refresh token secret not found")
	}

	return &tokenConfig{
		expTime: time.Minute * time.Duration(expTime),
		secret:  []byte(secret),
	}, nil
}

func NewAccessTokenConfig() (TokenConfig, error) {
	expTimeStr := os.Getenv(accessTokenExpTimeName)
	if expTimeStr == "" {
		return nil, errors.New("access token expiration time not found")
	}
	expTime, err := strconv.Atoi(expTimeStr)
	if err != nil {
		return nil, errors.New("unable to get access token expiration time")
	}

	secret := os.Getenv(accessTokenSecretKeyName)
	if secret == "" {
		return nil, errors.New("access token secret not found")
	}

	return &tokenConfig{
		expTime: time.Minute * time.Duration(expTime),
		secret:  []byte(secret),
	}, nil
}

func (cfg *tokenConfig) ExpTime() time.Duration {
	return cfg.expTime
}

func (cfg *tokenConfig) Secret() []byte {
	return cfg.secret
}
