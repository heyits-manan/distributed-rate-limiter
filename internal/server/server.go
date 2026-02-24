package server

import (
	"context"
	"net/http"

	"github.com/heyits-manan/distributed-rate-limiter/internal/config"
)

// Server wraps the HTTP server with graceful shutdown support.
type Server struct {
	httpServer *http.Server
}

// New creates a new Server with the given config and middleware chain.
func New(cfg config.ServerConfig, middlewares ...func(http.Handler) http.Handler) *Server {
	_ = http.NewServeMux() // TODO: assign to `mux` and register routes on it, e.g.:
	//   mux.HandleFunc("GET /healthz", healthHandler)

	// TODO: Wrap mux with your middleware chain
	// (hint: loop through middlewares and wrap: handler = mw(handler))

	// TODO: Create and return the Server with http.Server configured
	//   - Addr: fmt.Sprintf(":%d", cfg.Port)
	//   - Handler: your wrapped handler
	//   - ReadTimeout, WriteTimeout, IdleTimeout from cfg

	return nil
}

// Run starts the server and blocks until ctx is cancelled (e.g. Ctrl+C).
// It then shuts down the server gracefully.
func (s *Server) Run(ctx context.Context) error {
	// TODO:
	// 1. Start s.httpServer.ListenAndServe() in a goroutine
	// 2. Wait for either:
	//    a. The server to fail with an error → return the error
	//    b. ctx to be cancelled (Ctrl+C) → call s.httpServer.Shutdown() gracefully
	// (hint: use a channel and select{})
	return nil
}
