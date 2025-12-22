package http

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter gère le rate limiting par IP
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewRateLimiter crée un nouveau rate limiter
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate.Limit(float64(requestsPerMinute) / 60.0), // convertir en requêtes par seconde
		burst:    requestsPerMinute / 2,                         // burst = moitié de la limite
	}

	// Nettoyer les visiteurs inactifs toutes les minutes
	go rl.cleanupVisitors()

	return rl
}

// Allow vérifie si une requête depuis cette IP est autorisée
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
		return limiter.Allow()
	}

	v.lastSeen = time.Now()
	return v.limiter.Allow()
}

// cleanupVisitors supprime les visiteurs inactifs depuis plus de 3 minutes
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// rateLimitMiddleware applique le rate limiting
func rateLimitMiddleware(limiter *RateLimiter) func(next HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(w ResponseWriter, r *Request) error {
			ip := getClientIP(r.Request)

			if !limiter.Allow(ip) {
				w.WriteHeader(429)
				w.Write([]byte("Trop de requêtes. Veuillez réessayer plus tard."))
				return nil
			}

			return next(w, r)
		}
	}
}
