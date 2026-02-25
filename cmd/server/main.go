package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/heyits-manan/distributed-rate-limiter/internal/config"
	"github.com/heyits-manan/distributed-rate-limiter/internal/limiter"
	"github.com/heyits-manan/distributed-rate-limiter/internal/middleware"
	"github.com/heyits-manan/distributed-rate-limiter/internal/server"
	"github.com/heyits-manan/distributed-rate-limiter/internal/store"
)

func main() {
	// TODO (in this order):
	//
	// 1. Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	st := store.NewShardedStore(cfg.Store.ShardCount, cfg.Store.GCInterval)
	defer st.Close()

	rl := limiter.NewFixedWindow(st, cfg.RateLimit.RequestsPerWindow, cfg.RateLimit.WindowSize)

	mw := middleware.RateLimit(rl)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := server.New(cfg.Server, mw)
	if err := srv.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
