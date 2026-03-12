package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ipLimiter holds a rate limiter and the last time it was accessed.
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiterStore manages per-IP rate limiters with automatic cleanup.
type RateLimiterStore struct {
	mu       sync.Mutex
	limiters map[string]*ipLimiter
	rate     rate.Limit
	burst    int
}

// NewRateLimiterStore creates a new store with the given requests-per-second rate and burst size.
// It also starts a background goroutine to clean up stale entries every 3 minutes.
func NewRateLimiterStore(r rate.Limit, burst int) *RateLimiterStore {
	store := &RateLimiterStore{
		limiters: make(map[string]*ipLimiter),
		rate:     r,
		burst:    burst,
	}

	// Clean up stale limiters every 3 minutes
	go func() {
		for {
			time.Sleep(3 * time.Minute)
			store.mu.Lock()
			for ip, entry := range store.limiters {
				if time.Since(entry.lastSeen) > 5*time.Minute {
					delete(store.limiters, ip)
				}
			}
			store.mu.Unlock()
		}
	}()

	return store
}

// getLimiter returns the rate limiter for the given IP, creating one if it doesn't exist.
func (s *RateLimiterStore) getLimiter(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.limiters[ip]
	if !exists {
		limiter := rate.NewLimiter(s.rate, s.burst)
		s.limiters[ip] = &ipLimiter{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	entry.lastSeen = time.Now()
	return entry.limiter
}

// RateLimiter returns a Gin middleware that rate limits requests per client IP.
//
// Parameters:
//   - rps: requests per second allowed per IP
//   - burst: maximum burst size (extra requests allowed in a short window)
//
// Usage:
//
//	r.Use(middleware.RateLimiter(10, 20))          // 10 req/s, burst of 20
//	authGroup.Use(middleware.RateLimiter(3, 5))    // stricter for auth routes
func RateLimiter(rps float64, burst int) gin.HandlerFunc {
	store := NewRateLimiterStore(rate.Limit(rps), burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := store.getLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please slow down and try again shortly.",
			})
			return
		}

		c.Next()
	}
}
