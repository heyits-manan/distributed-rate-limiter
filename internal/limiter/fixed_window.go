package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/heyits-manan/distributed-rate-limiter/internal/store"
)

// FixedWindow implements RateLimiter using a fixed window algorithm.
// It counts requests in discrete time buckets (e.g. 0:00-1:00, 1:00-2:00).
type FixedWindow struct {
	store  store.Store
	limit  int
	window time.Duration
}

// NewFixedWindow creates a new fixed window rate limiter.
func NewFixedWindow(s store.Store, limit int, window time.Duration) *FixedWindow {
	return &FixedWindow{
		store:  s,
		limit:  limit,
		window: window,
	}
}

func (fw *FixedWindow) Allow(ctx context.Context, key string) (*Result, error) {
	now := time.Now()
	windowStart := now.Truncate(fw.window)
	windowKey := fmt.Sprintf("%s:%d", key, windowStart.Unix())
	resetAt := windowStart.Add(fw.window)
	count, error := fw.store.Increment(ctx, windowKey, fw.window)
	if error != nil {
		return nil, error
	}
	if count > fw.limit {
		return &Result{
			Allowed:    false,
			Limit:      fw.limit,
			Remaining:  0,
			RetryAfter: resetAt.Sub(now),
			ResetAt:    resetAt,
		}, nil
	}

	return &Result{
		Allowed:   true,
		Limit:     fw.limit,
		Remaining: fw.limit - count,
		ResetAt:   resetAt,
	}, nil
}
