package auth

import (
	"github.com/pkg/errors"

	"github.com/valek177/auth/grpc/pkg/auth_v1"
)

func validateLogin(req *auth_v1.LoginRequest) error {
	if req == nil {
		return errors.New("unable to login: empty request")
	}

	return nil
}

func validateRefreshTokenRequest(req *auth_v1.GetRefreshTokenRequest) error {
	if req == nil {
		return errors.New("unable to get refresh token: empty request")
	}

	return nil
}

func validateAccessTokenRequest(req *auth_v1.GetAccessTokenRequest) error {
	if req == nil {
		return errors.New("unable to get access token: empty request")
	}

	return nil
}
