package auth

import (
	"context"

	"github.com/valek177/auth/internal/model"
)

// UpdateUser updates user in repo
func (s *serv) UpdateUser(ctx context.Context, updateUserInfo *model.UpdateUserInfo) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.authRepository.UpdateUser(ctx, updateUserInfo)
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
