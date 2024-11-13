package auth

import (
	"context"

	"github.com/valek177/auth/grpc/pkg/auth_v1"
)

// Login returns refresh token for user
func (i *Implementation) Login(ctx context.Context, req *auth_v1.LoginRequest) (
	*auth_v1.LoginResponse, error,
) {
	err := validateLogin(req)
	if err != nil {
		return nil, err
	}

	refreshToken, err := i.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &auth_v1.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
