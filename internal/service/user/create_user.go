package user

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/valek177/auth/internal/converter"
	"github.com/valek177/auth/internal/model"
)

// CreateUser creates user in repo
func (s *serv) CreateUser(ctx context.Context, newUser *model.NewUser) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "create user (service)")
	defer span.Finish()

	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = validateCreateUser(newUser)
		if errTx != nil {
			return errTx
		}

		id, errTx = s.userRepository.CreateUser(ctx, newUser)
		if errTx != nil {
			return errTx
		}

		user, errTx := s.userRepository.GetUser(ctx, id)
		if errTx != nil {
			return errTx
		}

		_ = s.redisRepository.CreateUser(ctx, user)
		_ = s.redisRepository.SetExpireUser(ctx, id)

		_, errTx = s.logRepository.CreateRecord(ctx,
			converter.ToRecordRepoFromService(id, "create"))
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
