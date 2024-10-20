package auth

import (
	"context"
)

// DeleteUser deletes user in repo
func (s *serv) DeleteUser(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.authRepository.DeleteUser(ctx, id)
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
