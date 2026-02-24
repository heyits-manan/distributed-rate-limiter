package store

import (
	"context"
	"sync"
	"time"
)

// shard is one bucket of the sharded map.
// Each shard has its own lock so different keys don't block each other.
type shard struct {
	mu         sync.RWMutex
	counters   map[string]*counterEntry
	timestamps map[string][]time.Time
}

// counterEntry holds a count and its expiration time.
type counterEntry struct {
	count     int
	expiresAt time.Time
}

// ShardedStore splits keys across multiple shards to reduce lock contention.
type ShardedStore struct {
	shards []shard
	// TODO: Add fields for:
	// - a way to stop the background GC goroutine (hint: context.CancelFunc)
	// - a way to wait for the GC goroutine to finish (hint: sync.WaitGroup)
}

// NewShardedStore creates a new store with the given number of shards
// and starts a background GC goroutine.
func NewShardedStore(shardCount int, gcInterval time.Duration) *ShardedStore {
	// TODO: Create the shards slice, initialize each shard's maps
	// TODO: Start a background goroutine that calls sweep() every gcInterval
	// TODO: Return the store
	return nil
}

// getShard returns the shard that owns the given key.
func (s *ShardedStore) getShard(key string) *shard {
	// TODO: Hash the key (hint: hash/fnv) and pick a shard index
	return nil
}

func (s *ShardedStore) Increment(ctx context.Context, key string, expiration time.Duration) (int, error) {
	// TODO: Get the shard for this key
	// TODO: Lock the shard
	// TODO: If the key doesn't exist or is expired, create a new entry
	// TODO: Increment the counter
	// TODO: Unlock and return the count
	return 0, nil
}

func (s *ShardedStore) AddTimestamp(ctx context.Context, key string, t time.Time, window time.Duration) error {
	// TODO: Get the shard for this key
	// TODO: Lock the shard
	// TODO: Remove timestamps older than the window (pruning)
	// TODO: Append the new timestamp
	// TODO: Unlock and return
	return nil
}

func (s *ShardedStore) CountInWindow(ctx context.Context, key string, start, end time.Time) (int, error) {
	// TODO: Get the shard for this key
	// TODO: Read-lock the shard (hint: sh.mu.RLock() — allows parallel readers)
	// TODO: Count timestamps between start and end
	// TODO: Unlock and return the count
	return 0, nil
}

// Close stops the background GC and waits for it to finish.
func (s *ShardedStore) Close() error {
	// TODO: Cancel the GC goroutine's context
	// TODO: Wait for it to finish
	return nil
}

// sweep walks all shards and removes expired entries.
// Called periodically by the background GC goroutine.
func (s *ShardedStore) sweep() {
	// TODO: For each shard:
	//   - Lock it
	//   - Delete expired counters
	//   - Delete stale timestamp lists
	//   - Unlock it
}
