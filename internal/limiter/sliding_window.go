package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/itsmanan/distributed-rate-limiter/internal/config"
	"github.com/itsmanan/distributed-rate-limiter/internal/store"
)

type SlidingWindow struct {
	store  store.Store
	limit  int
	window time.Duration
}

func NewSlidingWindow(s store.Store, cfg config.RateLimitConfig) *SlidingWindow {
	return &SlidingWindow{
		store:  s,
		limit:  cfg.RequestsPerWindow,
		window: cfg.WindowSize,
	}
}

func (sw *SlidingWindow) Allow(ctx context.Context, key string) (*Result, error) {
	now := time.Now()
	windowStart := now.Add(-sw.window)
	storeKey := fmt.Sprintf("rl:sw:%s", key)

	count, err := sw.store.CountInWindow(ctx, storeKey, windowStart, now)
	if err != nil {
		return nil, fmt.Errorf("counting requests: %w", err)
	}

	result := &Result{
		Limit:   sw.limit,
		ResetAt: now.Add(sw.window),
	}

	if count >= sw.limit {
		result.Allowed = false
		result.Remaining = 0
		result.RetryAfter = sw.window - now.Sub(windowStart)
		return result, nil
	}

	if err := sw.store.AddTimestamp(ctx, storeKey, now, sw.window); err != nil {
		return nil, fmt.Errorf("recording request: %w", err)
	}

	result.Allowed = true
	result.Remaining = sw.limit - int(count) - 1

	return result, nil
}
