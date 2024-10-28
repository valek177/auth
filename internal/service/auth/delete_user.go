package auth

import (
	"context"

	"github.com/valek177/auth/internal/converter"
)

// DeleteUser deletes user in repo
func (s *serv) DeleteUser(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.authRepository.DeleteUser(ctx, id)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.logRepository.CreateRecord(ctx,
			converter.ToRecordRepoFromService(id, "delete"))
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
