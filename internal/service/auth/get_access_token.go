package auth

import (
	"context"
	"fmt"

	"github.com/valek177/auth/internal/model"
)

// GetAccessToken returns access token by refresh token
func (s *serv) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := s.tokenRefresh.VerifyToken(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %v", err)
	}

	accessToken, err := s.tokenAccess.GenerateToken(ctx, &model.User{
		Name: claims.Username,
		Role: claims.Role,
	})
	if err != nil {
		return "", fmt.Errorf("unable to generate access token: %v", err)
	}

	return accessToken, nil
}
