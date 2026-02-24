package middleware

import (
	"net/http"

	"github.com/heyits-manan/distributed-rate-limiter/internal/limiter"
)

// RateLimit returns HTTP middleware that rate-limits requests using the given limiter.
// Usage: server.Use(middleware.RateLimit(myLimiter))
func RateLimit(rl limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO:
			// 1. Extract the client key from the request (hint: IP address from r.RemoteAddr)
			// 2. Call rl.Allow(r.Context(), key)
			// 3. Set response headers:
			//    - X-RateLimit-Limit (the max)
			//    - X-RateLimit-Remaining (how many left)
			//    - X-RateLimit-Reset (when the window resets, as unix timestamp)
			// 4. If NOT allowed:
			//    - Set Retry-After header
			//    - Return 429 Too Many Requests
			// 5. If allowed:
			//    - Call next.ServeHTTP(w, r) to continue to the actual handler

			next.ServeHTTP(w, r)
		})
	}
}
