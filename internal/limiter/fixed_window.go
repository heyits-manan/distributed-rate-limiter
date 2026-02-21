package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/itsmanan/distributed-rate-limiter/internal/config"
	"github.com/itsmanan/distributed-rate-limiter/internal/store"
)

type FixedWindow struct {
	store  store.Store
	limit  int
	window time.Duration
}

func NewFixedWindow(s store.Store, cfg config.RateLimitConfig) *FixedWindow {
	return &FixedWindow{
		store:  s,
		limit:  cfg.RequestsPerWindow,
		window: cfg.WindowSize,
	}
}

func (fw *FixedWindow) Allow(ctx context.Context, key string) (*Result, error) {
	now := time.Now()
	windowKey := fmt.Sprintf("rl:fw:%s:%d", key, now.Truncate(fw.window).Unix())

	count, err := fw.store.Increment(ctx, windowKey, fw.window)
	if err != nil {
		return nil, fmt.Errorf("incrementing counter: %w", err)
	}

	result := &Result{
		Limit:   fw.limit,
		ResetAt: now.Truncate(fw.window).Add(fw.window),
	}

	if count > fw.limit {
		result.Allowed = false
		result.Remaining = 0
		result.RetryAfter = result.ResetAt.Sub(now)
		return result, nil
	}

	result.Allowed = true
	result.Remaining = fw.limit - count

	return result, nil
}
