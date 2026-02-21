package store

import (
	"context"
	"hash/fnv"
	"io"
	"time"
)

type Store interface {
	Increment(ctx context.Context, key string, expiration time.Duration) (int, error)

	AddTimestamp(ctx context.Context, key string, t time.Time, window time.Duration) error

	CountInWindow(ctx context.Context, key string, start, end time.Time) (int, error)

	io.Closer
}

type StoreConfig struct {
	ShardCount int
	GCInterval time.Duration
}

func DefaultStoreConfig() StoreConfig {
	return StoreConfig{
		ShardCount: 64,
		GCInterval: 30 * time.Second,
	}
}

func shardIndex(key string, shardCount int) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32()) % shardCount
}
