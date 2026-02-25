package store

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// --- Basic Tests ---

func TestIncrement(t *testing.T) {
	s := NewShardedStore(64, 30*time.Second)
	defer s.Close()

	ctx := context.Background()

	// First call should return 1
	count, err := s.Increment(ctx, "test-key", 60*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	// Second call should return 2
	count, err = s.Increment(ctx, "test-key", 60*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
}

func TestIncrementExpiration(t *testing.T) {
	s := NewShardedStore(64, 30*time.Second)
	defer s.Close()

	ctx := context.Background()

	// Increment with very short expiration
	s.Increment(ctx, "expire-key", 50*time.Millisecond)

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should reset to 1 (expired)
	count, _ := s.Increment(ctx, "expire-key", 50*time.Millisecond)
	if count != 1 {
		t.Errorf("expected count 1 after expiration, got %d", count)
	}
}

// --- Benchmarks ---

// BenchmarkIncrement_SingleKey — all goroutines hitting 1 key (worst case contention)
func BenchmarkIncrement_SingleKey(b *testing.B) {
	s := NewShardedStore(64, 30*time.Second)
	defer s.Close()

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Increment(ctx, "same-key", 60*time.Second)
		}
	})
}

// BenchmarkIncrement_UniqueKeys — each goroutine hits different keys (best case)
func BenchmarkIncrement_UniqueKeys(b *testing.B) {
	s := NewShardedStore(64, 30*time.Second)
	defer s.Close()

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", i)
			s.Increment(ctx, key, 60*time.Second)
			i++
		}
	})
}

// BenchmarkIncrement_ShardCount — compare different shard counts
func BenchmarkIncrement_ShardCount(b *testing.B) {
	for _, shards := range []int{1, 4, 16, 64, 256} {
		b.Run(fmt.Sprintf("shards-%d", shards), func(b *testing.B) {
			s := NewShardedStore(shards, 30*time.Second)
			defer s.Close()

			ctx := context.Background()

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					key := fmt.Sprintf("key-%d", i%100)
					s.Increment(ctx, key, 60*time.Second)
					i++
				}
			})
		})
	}
}
