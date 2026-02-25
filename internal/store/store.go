package store

import (
	"context"
	"io"
	"time"
)

// Store is the interface that all backends must implement.
type Store interface {
	// Increment adds 1 to a counter identified by key.
	// Returns the new count. The counter expires after `expiration`.
	Increment(ctx context.Context, key string, expiration time.Duration) (int, error)

	// Close cleans up resources (stop background goroutines, etc).
	io.Closer
}
