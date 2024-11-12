package user

import (
	"context"

	"github.com/pkg/errors"

	"github.com/valek177/auth/internal/model"
)

// GetUser returns user model from repo
func (s *serv) GetUser(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.redisRepository.GetUser(ctx, id)
	if err == nil {
		_ = s.redisRepository.SetExpireUser(ctx, id)
		return user, nil
	}

	user, err = s.userRepository.GetUser(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_ = s.redisRepository.CreateUser(ctx, user)
	_ = s.redisRepository.SetExpireUser(ctx, id)

	return user, nil
}
