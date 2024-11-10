package kafka

import (
	"context"

	"github.com/valek177/auth/internal/client/kafka/consumer"
)

// Consumer is interface for kafka consumer
type Consumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) (err error)
	Close() error
}
