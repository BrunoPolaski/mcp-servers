package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	httphelper "github.com/BrunoPolaski/registration-validation/internal/interfaces/http"
	"golang.org/x/time/rate"
)

// RateLimiter holds rate limiters per IP address
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
// rate: requests per second (e.g., 10 = 10 req/s)
// burst: maximum burst size (e.g., 20 = allow bursts up to 20 requests)
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

// getLimiter returns the rate limiter for the given IP
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
	}

	return limiter
}

// cleanup removes old limiters to prevent memory leaks
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for range ticker.C {
			rl.mu.Lock()
			// In production, implement proper cleanup logic
			// For now, just clear if too many entries
			if len(rl.limiters) > 10000 {
				rl.limiters = make(map[string]*rate.Limiter)
			}
			rl.mu.Unlock()
		}
	}()
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(limiter *RateLimiter) Middleware {
	limiter.cleanup()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP from various headers (for proxy support)
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = r.Header.Get("X-Real-IP")
			}
			if ip == "" {
				ip = r.RemoteAddr
			}

			limiter := limiter.getLimiter(ip)
			if !limiter.Allow() {
				httphelper.ErrorResponse(
					rest_err.NewRestErr("rate limit exceeded, please try again later", "too_many_requests", http.StatusTooManyRequests, nil),
					w,
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
