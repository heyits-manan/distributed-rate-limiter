package store

import (
	"context"
	"log/slog"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type counterEntry struct {
	count     int
	expiresAt time.Time
}

type shard struct {
	mu         sync.RWMutex
	counters   map[string]*counterEntry
	timestamps map[string][]time.Time
}

type ShardedStore struct {
	shards    []shard
	shardMask int
	cancel    context.CancelFunc
	wg        sync.WaitGroup

	hits   atomic.Int64
	misses atomic.Int64
	evicts atomic.Int64
}

func NewShardedStore(cfg StoreConfig) *ShardedStore {
	n := nextPowerOfTwo(cfg.ShardCount)

	shards := make([]shard, n)
	for i := range shards {
		shards[i] = shard{
			counters:   make(map[string]*counterEntry),
			timestamps: make(map[string][]time.Time),
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	s := &ShardedStore{
		shards:    shards,
		shardMask: n - 1,
		cancel:    cancel,
	}

	s.wg.Add(1)
	go s.gcLoop(ctx, cfg.GCInterval)

	return s
}

func (s *ShardedStore) getShard(key string) *shard {
	idx := shardIndex(key, len(s.shards))
	return &s.shards[idx&s.shardMask]
}

// --- Store interface implementation ---

func (s *ShardedStore) Increment(_ context.Context, key string, expiration time.Duration) (int, error) {
	sh := s.getShard(key)

	sh.mu.Lock()
	defer sh.mu.Unlock()

	now := time.Now()

	e, ok := sh.counters[key]
	if !ok || now.After(e.expiresAt) {
		if ok {
			s.evicts.Add(1)
		}
		s.misses.Add(1)
		e = &counterEntry{expiresAt: now.Add(expiration)}
		sh.counters[key] = e
	} else {
		s.hits.Add(1)
	}

	e.count++
	return e.count, nil
}

func (s *ShardedStore) AddTimestamp(_ context.Context, key string, t time.Time, window time.Duration) error {
	sh := s.getShard(key)

	sh.mu.Lock()
	defer sh.mu.Unlock()

	cutoff := t.Add(-window)
	existing := sh.timestamps[key]

	pruneIdx := sort.Search(len(existing), func(i int) bool {
		return existing[i].After(cutoff)
	})

	pruned := existing[pruneIdx:]
	evicted := int64(pruneIdx)
	if evicted > 0 {
		s.evicts.Add(evicted)
	}

	sh.timestamps[key] = append(pruned, t)

	return nil
}

func (s *ShardedStore) CountInWindow(_ context.Context, key string, start, end time.Time) (int, error) {
	sh := s.getShard(key)

	sh.mu.RLock()
	defer sh.mu.RUnlock()

	ts := sh.timestamps[key]
	if len(ts) == 0 {
		s.misses.Add(1)
		return 0, nil
	}

	s.hits.Add(1)

	lo := sort.Search(len(ts), func(i int) bool {
		return !ts[i].Before(start)
	})
	hi := sort.Search(len(ts), func(i int) bool {
		return ts[i].After(end)
	})

	return hi - lo, nil
}

func (s *ShardedStore) Close() error {
	s.cancel()
	s.wg.Wait()
	return nil
}

// --- Background GC ---

func (s *ShardedStore) gcLoop(ctx context.Context, interval time.Duration) {
	defer s.wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.sweep()
		}
	}
}

func (s *ShardedStore) sweep() {
	now := time.Now()
	var totalEvicted int64

	for i := range s.shards {
		sh := &s.shards[i]
		sh.mu.Lock()

		for key, entry := range sh.counters {
			if now.After(entry.expiresAt) {
				delete(sh.counters, key)
				totalEvicted++
			}
		}

		for key, ts := range sh.timestamps {
			if len(ts) == 0 {
				delete(sh.timestamps, key)
				continue
			}
			if now.Sub(ts[len(ts)-1]) > 5*time.Minute {
				delete(sh.timestamps, key)
				totalEvicted += int64(len(ts))
			}
		}

		sh.mu.Unlock()
	}

	if totalEvicted > 0 {
		s.evicts.Add(totalEvicted)
		slog.Debug("gc sweep completed", "evicted", totalEvicted)
	}
}

// --- Metrics ---

type Metrics struct {
	Hits       int64
	Misses     int64
	Evictions  int64
	ShardCount int
	Keys       int
}

func (s *ShardedStore) Metrics() Metrics {
	var totalKeys int
	for i := range s.shards {
		sh := &s.shards[i]
		sh.mu.RLock()
		totalKeys += len(sh.counters) + len(sh.timestamps)
		sh.mu.RUnlock()
	}

	return Metrics{
		Hits:       s.hits.Load(),
		Misses:     s.misses.Load(),
		Evictions:  s.evicts.Load(),
		ShardCount: len(s.shards),
		Keys:       totalKeys,
	}
}

// --- Helpers ---

func nextPowerOfTwo(n int) int {
	if n <= 0 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	return n + 1
}
