package limiter

import (
	"context"
	"time"

	"github.com/heyits-manan/distributed-rate-limiter/internal/store"
)

// SlidingWindow implements RateLimiter using a sliding window algorithm.
// It counts requests in a rolling time window (e.g. last 60 seconds).
type SlidingWindow struct {
	store  store.Store
	limit  int
	window time.Duration
}

// NewSlidingWindow creates a new sliding window rate limiter.
func NewSlidingWindow(s store.Store, limit int, window time.Duration) *SlidingWindow {
	return &SlidingWindow{
		store:  s,
		limit:  limit,
		window: window,
	}
}

func (sw *SlidingWindow) Allow(ctx context.Context, key string) (*Result, error) {
	// TODO:
	// 1. Get current time
	// 2. Calculate window start (now - window duration)
	// 3. Call sw.store.CountInWindow(ctx, key, windowStart, now)
	// 4. If count >= sw.limit → return denied result
	// 5. Call sw.store.AddTimestamp(ctx, key, now, sw.window) to record this request
	// 6. Return allowed result with remaining count
	return nil, nil
}
