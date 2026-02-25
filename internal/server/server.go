package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/heyits-manan/distributed-rate-limiter/internal/config"
)

// Server wraps the HTTP server with graceful shutdown support.
type Server struct {
	httpServer *http.Server
}

// New creates a new Server with the given config and middleware chain.
func New(cfg config.ServerConfig, middlewares ...func(http.Handler) http.Handler) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	var handler http.Handler = mux
	for _, mw := range middlewares {
		handler = mw(handler)
	}

	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
}

// Run starts the server and blocks until ctx is cancelled (e.g. Ctrl+C).
// It then shuts down the server gracefully.
func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return s.httpServer.Shutdown(context.Background())
	}
}
