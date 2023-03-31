package server

import (
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter interface
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// Return new IPRateLimiter .
func NewIPRateLimiter(r int, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   rate.Limit(r),
		b:   b,
	}

	return i
}

// AddIP creates a new rate limiter and adds it to the ips map,
// using the IP address as the key
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)

	i.ips[ip] = limiter

	return limiter
}

// GetLimiter returns the rate limiter for the provided IP address if it exists.
// Otherwise calls AddIP to add IP address to the map
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.ips[ip]

	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}

	i.mu.Unlock()

	return limiter
}

func limitMiddleware(limiter *IPRateLimiter) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		limiter := limiter.GetLimiter(ctx.ClientIP())
		if !limiter.Allow() {
			ctx.JSON(429, gin.H{})
			return
		}
	}
}
