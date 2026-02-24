package limiter

import (
	"context"
	"time"
)

// Result is returned by Allow() — tells the caller whether the request is allowed.
type Result struct {
	Allowed    bool          // true = request can proceed
	Limit      int           // max requests per window
	Remaining  int           // how many requests are left
	RetryAfter time.Duration // how long to wait if denied
	ResetAt    time.Time     // when the window resets
}

// RateLimiter is the interface that all rate limiting algorithms must implement.
type RateLimiter interface {
	// Allow checks if the given key (e.g. IP address) is allowed to make a request.
	Allow(ctx context.Context, key string) (*Result, error)
}
