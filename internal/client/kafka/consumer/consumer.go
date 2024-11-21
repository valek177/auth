package consumer

import (
	"context"
	"strings"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"

	"github.com/valek177/auth/internal/logger"
)

type consumer struct {
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *GroupHandler
}

// NewConsumer returns new kafka consumer
func NewConsumer(
	consumerGroup sarama.ConsumerGroup,
	consumerGroupHandler *GroupHandler,
) *consumer {
	return &consumer{
		consumerGroup:        consumerGroup,
		consumerGroupHandler: consumerGroupHandler,
	}
}

// Consume executes consume logic
func (c *consumer) Consume(ctx context.Context, topicName string, handler Handler) error {
	c.consumerGroupHandler.msgHandler = handler

	return c.consume(ctx, topicName)
}

// Close closes consumer group
func (c *consumer) Close() error {
	return c.consumerGroup.Close()
}

func (c *consumer) consume(ctx context.Context, topicName string) error {
	for {
		err := c.consumerGroup.Consume(ctx, strings.Split(topicName, ","), c.consumerGroupHandler)
		if err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		logger.Debug("rebalancing...")
	}
}
