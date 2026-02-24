package limiter

import (
	"context"
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
	// TODO:
	// 1. Get current time
	// 2. Build a window key like "rl:fw:<key>:<window_start_unix>"
	//    (hint: now.Truncate(fw.window) gives you the start of the current window)
	// 3. Call fw.store.Increment(ctx, windowKey, fw.window)
	// 4. If count > fw.limit → return denied result
	// 5. Return allowed result with remaining count
	return nil, nil
}
