package store

import (
	"context"
	"io"
	"time"
)

// Store is the interface that all backends must implement.
// The in-memory store and (later) Redis store both satisfy this.
type Store interface {
	// Increment adds 1 to a counter identified by key.
	// Returns the new count. The counter expires after `expiration`.
	// Used by: fixed window algorithm.
	Increment(ctx context.Context, key string, expiration time.Duration) (int, error)

	// AddTimestamp records that a request happened at time t.
	// Old timestamps outside the window should be pruned.
	// Used by: sliding window algorithm.
	AddTimestamp(ctx context.Context, key string, t time.Time, window time.Duration) error

	// CountInWindow returns how many requests happened between start and end.
	// Used by: sliding window algorithm.
	CountInWindow(ctx context.Context, key string, start, end time.Time) (int, error)

	// Close cleans up resources (stop background goroutines, etc).
	io.Closer
}
