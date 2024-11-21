package user_saver

import (
	"context"

	"github.com/valek177/auth/internal/client/kafka"
	"github.com/valek177/auth/internal/repository"
	def "github.com/valek177/auth/internal/service"
)

const (
	topicName = "test-topic"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
	consumer       kafka.Consumer
}

// NewService returns new consumer service
func NewService(
	userRepository repository.UserRepository,
	consumer kafka.Consumer,
) *service {
	return &service{
		userRepository: userRepository,
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

		errChan <- s.consumer.Consume(ctx, topicName, s.UserSaveHandler)
	}()

	return errChan
}
