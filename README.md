# Distributed Rate Limiter

A high-performance, distributed rate limiter built in Go. Supports multiple algorithms (sliding window, fixed window) with pluggable backends (Redis, in-memory).

## Features

- **Multiple algorithms** — sliding window and fixed window counters
- **Pluggable stores** — Redis for distributed deployments, in-memory for local dev
- **HTTP middleware** — drop-in `net/http` middleware
- **Graceful shutdown** — clean server lifecycle management
- **Docker-ready** — multi-stage build with distroless runtime

## Quick Start

```bash
# Run locally
make run

# Run with Docker Compose (app + Redis)
make docker-up

# Run tests
make test
```

## Project Structure

```
cmd/server/          → Application entrypoint
internal/config/     → Configuration loading
internal/limiter/    → Rate limiting algorithms
internal/store/      → Backend store implementations
internal/middleware/  → HTTP middleware
internal/server/     → HTTP server setup
pkg/ratelimit/       → Public client SDK
api/proto/           → Protobuf definitions
configs/             → Default config files
deployments/         → Docker Compose, k8s manifests
```

## Configuration

Configure via `configs/config.yaml` or environment variables:

| Variable          | Default        | Description              |
|-------------------|----------------|--------------------------|
| `SERVER_PORT`     | `8080`         | HTTP listen port         |
| `REDIS_ADDR`      | `localhost:6379` | Redis address          |
| `RATE_LIMIT`      | `100`          | Requests per window      |
| `WINDOW_SIZE`     | `60s`          | Window duration          |

## License

MIT
