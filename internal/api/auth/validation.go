package auth

import (
	"github.com/pkg/errors"

	"github.com/valek177/auth/grpc/pkg/user_v1"
)

func validateCreateUser(req *user_v1.CreateUserRequest) error {
	if req == nil {
		return errors.New("unable to create user: empty request")
	}

	return nil
}

func validateDeleteUser(req *user_v1.DeleteUserRequest) error {
	if req == nil {
		return errors.New("unable to delete user: empty request")
	}

	return nil
}

func validateUpdateUser(req *user_v1.UpdateUserRequest) error {
	if req == nil {
		return errors.New("unable to update user: empty request")
	}

	return nil
}

func validateGetUser(req *user_v1.GetUserRequest) error {
	if req == nil {
		return errors.New("unable to get user: empty request")
	}

	return nil
}
