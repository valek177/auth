package user_saver

import (
	"context"

	"github.com/valek177/auth/internal/client/kafka"
	"github.com/valek177/auth/internal/repository"
	def "github.com/valek177/auth/internal/service"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	authRepository repository.AuthRepository
	consumer       kafka.Consumer
}

// NewService returns new consumer service
func NewService(
	authRepository repository.AuthRepository,
	consumer kafka.Consumer,
) *service {
	return &service{
		authRepository: authRepository,
		consumer:       consumer,
	}
}

// RunConsumer runs consumer logic
func (s *service) RunConsumer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-s.run(ctx):
			if err != nil {
				return err
			}
		}
	}
}

func (s *service) run(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		errChan <- s.consumer.Consume(ctx, "test-topic", s.UserSaveHandler)
	}()

	return errChan
}
