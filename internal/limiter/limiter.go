package limiter

import (
	"context"
	"time"
)

type Result struct {
	Allowed    bool
	Limit      int
	Remaining  int
	RetryAfter time.Duration
	ResetAt    time.Time
}

type RateLimiter interface {
	Allow(ctx context.Context, key string) (*Result, error)
}
