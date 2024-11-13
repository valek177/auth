package auth

import (
	"context"

	"github.com/valek177/auth/grpc/pkg/auth_v1"
)

func (i *Implementation) GetRefreshToken(ctx context.Context, req *auth_v1.GetRefreshTokenRequest) (
	*auth_v1.GetRefreshTokenResponse, error,
) {
	err := validateRefreshTokenRequest(req)
	if err != nil {
		return nil, err
	}

	refreshToken, err := i.authService.GetRefreshToken(ctx, req.GetOldRefreshToken())
	if err != nil {
		return nil, err
	}

	return &auth_v1.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}
