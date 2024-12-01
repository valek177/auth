package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rateLimiter "github.com/valek177/auth/internal/rate_limiter"
)

// RateLimiterInterceptor is a struct for rate limiter interceptor
type RateLimiterInterceptor struct {
	rateLimiter *rateLimiter.TokenBucketLimiter
}

// NewRateLimiterInterceptor creates new rate limiter interceptor
func NewRateLimiterInterceptor(rateLimiter *rateLimiter.TokenBucketLimiter) *RateLimiterInterceptor {
	return &RateLimiterInterceptor{rateLimiter: rateLimiter}
}

// Unary returns unary interface
func (r *RateLimiterInterceptor) Unary(ctx context.Context, req interface{},
	_ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	if !r.rateLimiter.Allow() {
		return nil, status.Error(codes.ResourceExhausted, "too many requests")
	}

	return handler(ctx, req)
}
