package auth

import (
	"context"
	"fmt"

	passwordLib "github.com/valek177/auth/internal/password"
)

// Login returns refresh and access tokens for username & password
func (s *serv) Login(ctx context.Context, username, password string) (string, string, error) {
	user, err := s.userRepository.GetUserByName(ctx, username)
	if err != nil {
		return "", "", fmt.Errorf("unable to get user %s", username)
	}

	isPasswordsEqual := passwordLib.CheckPasswordHash(password, user.Password)
	if !isPasswordsEqual {
		return "", "", fmt.Errorf("unable to login: incorrect password")
	}

	refreshToken, err := s.tokenRefresh.GenerateToken(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token")
	}

	accessToken, err := s.tokenAccess.GenerateToken(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token")
	}

	return refreshToken, accessToken, nil
}
