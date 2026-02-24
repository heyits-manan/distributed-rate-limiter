package main

import (
	"fmt"
	"os"

	"github.com/heyits-manan/distributed-rate-limiter/internal/config"
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
	_ = cfg // remove this line once you use cfg

	// 2. Create the store
	//    st := store.NewShardedStore(cfg.Store.ShardCount, cfg.Store.GCInterval)
	//    defer st.Close()

	// 3. Create the limiter
	//    rl := limiter.NewSlidingWindow(st, cfg.RateLimit.RequestsPerWindow, cfg.RateLimit.WindowSize)

	// 4. Create the middleware
	//    mw := middleware.RateLimit(rl)

	// 5. Create and run the server
	//    srv := server.New(cfg.Server, mw)
	//    srv.Run(ctx)
}
