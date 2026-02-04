package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/invopop/ctxi18n/i18n"
	"github.com/rousseau-romain/round-timing/handlers"
)

type visitor struct {
	count    int
	lastSeen time.Time
}

// RateLimiter tracks request counts per IP within a sliding window.
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	max      int
	window   time.Duration
}

// NewRateLimiter creates a rate limiter allowing max requests per window per IP.
func NewRateLimiter(max int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		max:      max,
		window:   window,
	}
	go rl.cleanup()
	return rl
}

// cleanup removes stale entries every minute.
func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Limit wraps a handler and rejects requests that exceed the rate limit.
func (rl *RateLimiter) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		if !exists || time.Since(v.lastSeen) > rl.window {
			rl.visitors[ip] = &visitor{count: 1, lastSeen: time.Now()}
			rl.mu.Unlock()
			next(w, r)
			return
		}

		v.count++
		v.lastSeen = time.Now()

		if v.count > rl.max {
			rl.mu.Unlock()
			errMessage := i18n.T(r.Context(), "page.signin.too-many-requests")
			handlers.RenderComponentError(errMessage, []string{errMessage}, http.StatusTooManyRequests, w, r)
			return
		}
		rl.mu.Unlock()

		next(w, r)
	}
}
