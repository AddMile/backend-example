package ratelimit

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const cleanupTTL = 5 * time.Minute

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
}

func New() *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
	}

	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)

		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > cleanupTTL {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) AllowedByIP(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, ok := rl.visitors[ip]
	if !ok {
		// 10 req per second with max burst 20
		limiter := rate.NewLimiter(100, 200)
		rl.visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}

		return true
	}
	v.lastSeen = time.Now()

	return v.limiter.Allow()
}
