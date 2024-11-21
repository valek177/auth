package auth

import (
	"context"

	"github.com/valek177/auth/grpc/pkg/auth_v1"
)

// GetAccessToken returns access token by refresh token
func (i *Implementation) GetAccessToken(ctx context.Context, req *auth_v1.GetAccessTokenRequest) (
	*auth_v1.GetAccessTokenResponse, error,
) {
	err := validateAccessTokenRequest(req)
	if err != nil {
		return nil, err
	}

	accessToken, err := i.authService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &auth_v1.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
