package user

import (
	"github.com/pkg/errors"

	"github.com/valek177/auth/internal/model"
	"github.com/valek177/auth/internal/password"
)

func validateCreateUser(user *model.NewUser) error {
	if user == nil {
		return errors.New("unable to create user: empty model")
	}
	if user.Name == "" {
		return errors.New("unable to create user: name is required")
	}
	if user.Password == "" {
		return errors.New("unable to create user: password is required")
	}
	if !password.CheckPasswordHash(user.PasswordConfirm, user.Password) {
		return errors.New("unable to create user: the passwords do not match")
	}

	return nil
}

func validateUpdateUser(user *model.UpdateUserInfo) error {
	if user == nil {
		return errors.New("unable to update user: empty model")
	}
	if user.Name == nil && user.Role == nil {
		return errors.New("unable to update user: nothing to update")
	}

	return nil
}
