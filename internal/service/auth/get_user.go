package auth

import (
	"context"

	"github.com/valek177/auth/internal/model"
)

// GetUser returns user model from repo
func (s *serv) GetUser(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.authRepository.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
