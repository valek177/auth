package user_saver

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/valek177/auth/internal/logger"
	"github.com/valek177/auth/internal/model"
)

// UserSaveHandler executes kafka creation user logic
func (s *service) UserSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	user := &model.NewUser{}
	err := json.Unmarshal(msg.Value, user)
	if err != nil {
		return err
	}

	id, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	logger.Debug("Kafka user handler: user with id was created", zap.Int64("id", id))

	return nil
}
