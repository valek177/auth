package auth

import (
	"context"

	"github.com/valek177/auth/internal/model"
)

// CreateUser creates user in repo
func (s *serv) CreateUser(ctx context.Context, newUser *model.NewUser) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.authRepository.CreateUser(ctx, newUser)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.authRepository.GetUser(ctx, id)
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}
