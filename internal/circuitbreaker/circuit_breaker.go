package circuitbreaker

import (
	"time"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"

	"github.com/valek177/auth/internal/logger"
)

func NewCircuitBreaker() *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "my-service",
		MaxRequests: 3,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Debug("Circuit Breaker, changed state",
				zap.String("name", name), zap.String("from", from.String()),
				zap.String("to", to.String()))
		},
	})
}
