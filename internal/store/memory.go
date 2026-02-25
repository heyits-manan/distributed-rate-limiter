package store

import (
	"context"
	"hash/fnv"
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
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewShardedStore creates a new store with the given number of shards
// and starts a background GC goroutine.
func NewShardedStore(shardCount int, gcInterval time.Duration) *ShardedStore {
	shards := make([]shard, shardCount)
	for i := range shards {
		shards[i].counters = make(map[string]*counterEntry)
	}
	ctx, cancel := context.WithCancel(context.Background())

	s := &ShardedStore{
		shards: shards,
		cancel: cancel,
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		<-ctx.Done()
	}()

	return s
}

// getShard returns the shard that owns the given key.
func (s *ShardedStore) getShard(key string) *shard {
	h := fnv.New32a()
	h.Write([]byte(key))
	number := h.Sum32()
	index := number % uint32(len(s.shards))
	return &s.shards[index]
}

func (s *ShardedStore) Increment(ctx context.Context, key string, expiration time.Duration) (int, error) {
	sh := s.getShard(key)
	sh.mu.Lock()
	defer sh.mu.Unlock()

	entry, exists := sh.counters[key]
	if !exists || time.Now().After(entry.expiresAt) {
		entry = &counterEntry{
			count:     0,
			expiresAt: time.Now().Add(expiration),
		}
		sh.counters[key] = entry
	}
	entry.count++
	return entry.count, nil
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
	s.cancel()
	s.wg.Wait()
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
