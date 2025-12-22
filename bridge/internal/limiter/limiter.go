package limiter

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

type PerServiceLimiter struct {
	limiter *rate.Limiter
}

func NewLimiter(qps int, burst int) *PerServiceLimiter {
	return &PerServiceLimiter{
		limiter: rate.NewLimiter(rate.Limit(qps), burst),
	}
}

func (l *PerServiceLimiter) Allow() bool {
	return l.limiter.Allow()
}

func (l *PerServiceLimiter) Wait() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return l.limiter.Wait(ctx)
}
