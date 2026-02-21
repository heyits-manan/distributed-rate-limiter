# Architecture

## Overview

The distributed rate limiter is a horizontally scalable service that enforces request rate limits across multiple instances using a shared backend store (Redis).

## Components

```
┌─────────┐     ┌────────────┐     ┌──────────┐     ┌───────┐
│  Client  │────▶│ Middleware  │────▶│ Limiter  │────▶│ Store │
└─────────┘     └────────────┘     └──────────┘     └───────┘
                                        │                │
                                   ┌────┴────┐     ┌─────┴────┐
                                   │ Sliding  │     │  Redis   │
                                   │ Window   │     │  Memory  │
                                   │ Fixed    │     └──────────┘
                                   │ Window   │
                                   └──────────┘
```

### Middleware
Extracts the client key (IP address) from incoming HTTP requests and delegates to the limiter. Sets standard rate-limit response headers.

### Limiter
Implements rate limiting algorithms behind a `RateLimiter` interface. Current implementations:
- **Sliding Window** — counts requests in a rolling time window using sorted sets
- **Fixed Window** — counts requests in discrete time buckets using atomic counters

### Store
Abstracts the persistence layer behind a `Store` interface:
- **Redis** — production backend, enables distributed rate limiting across multiple instances
- **Memory** — for local development and testing

## Request Flow

1. HTTP request arrives at the server
2. Rate limit middleware extracts the client key
3. Middleware calls `limiter.Allow(ctx, key)`
4. Limiter queries the store for the current count
5. If under limit: request recorded, forwarded to handler
6. If over limit: `429 Too Many Requests` returned with `Retry-After` header
