# Distributed Rate Limiter

A high-performance, distributed rate limiter built in pure Go. Designed from the ground up to learn and apply concurrency control, lock granularity, memory lifecycle management, and performance profiling — no Redis, no shortcuts.

## Why Pure Go?

This project intentionally avoids external stores like Redis to force direct engagement with:

- **Concurrency control** — `sync.Mutex`, `sync.RWMutex`, and `atomic` operations
- **Lock granularity** — sharded maps with per-shard locks to reduce contention
- **Memory lifecycle** — background GC goroutine for expired entry cleanup
- **Performance profiling** — built-in `pprof` endpoints for mutex and heap analysis
- **Failure behavior** — graceful degradation and clean shutdown coordination

## Features

- **Multiple algorithms** — sliding window (sorted timestamps) and fixed window (atomic counters)
- **Sharded concurrent store** — 64 shards with per-shard `RWMutex` for fine-grained locking
- **Background GC** — configurable sweep interval for expired entries
- **HTTP middleware** — drop-in `net/http` middleware with standard rate-limit headers
- **Profiling ready** — `pprof` server on `:6060` out of the box
- **Graceful shutdown** — context-based lifecycle with `WaitGroup` coordination
- **Docker-ready** — multi-stage build with distroless runtime

## Quick Start

```bash
# Run locally
make run

# Run with Docker
make docker-up

# Run tests
make test

# Profile lock contention
go tool pprof http://localhost:6060/debug/pprof/mutex
```

## Project Structure

```
cmd/server/          → Application entrypoint
internal/config/     → Configuration loading
internal/limiter/    → Rate limiting algorithms (sliding window, fixed window)
internal/store/      → Sharded concurrent in-memory store
internal/middleware/  → HTTP rate-limit middleware
internal/server/     → HTTP server with graceful shutdown
pkg/ratelimit/       → Public client SDK
api/proto/           → Protobuf definitions (future)
configs/             → Default config files
deployments/         → Docker Compose
```

## Configuration

Configure via environment variables or `configs/config.yaml`:

| Variable             | Default | Description                        |
|----------------------|---------|------------------------------------|
| `SERVER_PORT`        | `8080`  | HTTP listen port                   |
| `STORE_SHARD_COUNT`  | `64`    | Number of map shards               |
| `STORE_GC_INTERVAL`  | `30s`   | Background GC sweep interval       |
| `RATE_LIMIT`         | `100`   | Requests per window                |
| `WINDOW_SIZE`        | `60s`   | Window duration                    |


## License

MIT
