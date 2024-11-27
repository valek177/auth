package user

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/valek177/auth/internal/converter"
)

// DeleteUser deletes user in repo
func (s *serv) DeleteUser(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "delete user (service)")
	defer span.Finish()

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.userRepository.DeleteUser(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.redisRepository.DeleteUser(ctx, id)
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
