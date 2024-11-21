package user

import (
	"context"

	"github.com/valek177/auth/internal/converter"
	"github.com/valek177/auth/internal/model"
)

// UpdateUser updates user in repo
func (s *serv) UpdateUser(ctx context.Context, updateUserInfo *model.UpdateUserInfo) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		errTx = validateUpdateUser(updateUserInfo)
		if errTx != nil {
			return errTx
		}

		errTx = s.userRepository.UpdateUser(ctx, updateUserInfo)
		if errTx != nil {
			return errTx
		}

		user, errTx := s.userRepository.GetUser(ctx, updateUserInfo.ID)
		if errTx != nil {
			return errTx
		}

		errTx = s.redisRepository.DeleteUser(ctx, user.ID)
		if errTx != nil {
			return errTx
		}
		_ = s.redisRepository.CreateUser(ctx, user)
		_ = s.redisRepository.SetExpireUser(ctx, user.ID)

		_, errTx = s.logRepository.CreateRecord(ctx,
			converter.ToRecordRepoFromService(updateUserInfo.ID, "update"))
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
