package middlewares

import (
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	clientLimits = make(map[string]*rate.Limiter)
	mu           sync.Mutex
	rateLimit    = rate.Limit(1) // 1 request per second
	burstLimit   = 5             // Allow bursts of up to 5 requests
)

// getClientLimiter retrieves the rate limiter for a specific client IP
func getClientLimiter(clientIP string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := clientLimits[clientIP]
	if !exists {
		limiter = rate.NewLimiter(rateLimit, burstLimit)
		clientLimits[clientIP] = limiter
	}
	return limiter
}

// RateLimitMiddleware applies rate limiting to the requests
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter := getClientLimiter(clientIP)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(429, gin.H{"error": "too many requests"})
			return
		}

		c.Next()
	}
}
