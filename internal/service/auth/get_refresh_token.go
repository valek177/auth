package auth

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/valek177/auth/internal/model"
)

// GetRefreshToken returns new refresh token by old refresh token
func (s *serv) GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get refresh token (service)")
	defer span.Finish()

	claims, err := s.tokenRefresh.VerifyToken(ctx, oldRefreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	refreshToken, err := s.tokenRefresh.GenerateToken(ctx, &model.User{
		Name: claims.Username,
		Role: claims.Role,
	})
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
