package main

import (
	"context"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsmanan/distributed-rate-limiter/internal/config"
	"github.com/itsmanan/distributed-rate-limiter/internal/limiter"
	"github.com/itsmanan/distributed-rate-limiter/internal/middleware"
	"github.com/itsmanan/distributed-rate-limiter/internal/server"
	"github.com/itsmanan/distributed-rate-limiter/internal/store"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	st := store.NewShardedStore(store.StoreConfig{
		ShardCount: cfg.Store.ShardCount,
		GCInterval: cfg.Store.GCInterval,
	})
	defer st.Close()

	go func() {
		slog.Info("pprof server starting", "addr", ":6060")
		if err := http.ListenAndServe(":6060", nil); err != nil {
			slog.Error("pprof server failed", "error", err)
		}
	}()

	rl := limiter.NewSlidingWindow(st, cfg.RateLimit)

	mw := middleware.RateLimit(rl)

	srv := server.New(cfg.Server, mw)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := srv.Run(ctx); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
