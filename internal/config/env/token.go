package env

import (
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	refreshTokenExpTimeName   = "REFRESH_TOKEN_EXPIRATION_TIME" //nolint:gosec
	refreshTokenSecretKeyName = "REFRESH_TOKEN_SECRET_KEY"      //nolint:gosec

	accessTokenExpTimeName   = "ACCESS_TOKEN_EXPIRATION_TIME" //nolint:gosec
	accessTokenSecretKeyName = "ACCESS_TOKEN_SECRET_KEY"      //nolint:gosec
)

// TokenConfig is interface for token config
type TokenConfig interface {
	ExpTime() time.Duration
	Secret() []byte
}

type tokenConfig struct {
	expTime time.Duration
	secret  []byte
}

// NewRefreshTokenConfig returns token config for refresh token
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

// NewAccessTokenConfig returns token config for access token
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

// ExpTime return expiration time for token
func (cfg *tokenConfig) ExpTime() time.Duration {
	return cfg.expTime
}

// Secret return secret for token
func (cfg *tokenConfig) Secret() []byte {
	return cfg.secret
}
