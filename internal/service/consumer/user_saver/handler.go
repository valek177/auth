package user_saver

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"

	"github.com/valek177/auth/internal/model"
)

// UserSaveHandler executes kafka creation user logic
func (s *service) UserSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	user := &model.NewUser{}
	err := json.Unmarshal(msg.Value, user)
	if err != nil {
		return err
	}

	id, err := s.authRepository.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	log.Printf("Kafka user handler: user with id %d was created\n", id)

	return nil
}
