package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alexisgalvan/onepiece-api/internal/domain"
	"github.com/alexisgalvan/onepiece-api/pkg/logger"
)

const (
	rateLimitPublicWindow = 1 * time.Minute
	rateLimitAuthWindow   = 1 * time.Minute
)

// RateLimitPublic limits requests per IP for public endpoints.
func RateLimitPublic(cache domain.CacheRepository, limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rl:public:%s", ip)

		allowed, err := cache.RateLimit(c.Request.Context(), key, limit, rateLimitPublicWindow)
		if err != nil {
			logger.Log.Error("rate limit check failed", zap.Error(err))
			c.Next() // fail open — don't block on Redis error
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))

		if !allowed {
			errMsg := "rate limit exceeded — try again in 1 minute"
			c.AbortWithStatusJSON(http.StatusTooManyRequests, domain.Response{
				Success: false,
				Data:    nil,
				Error:   &errMsg,
			})
			return
		}

		c.Next()
	}
}

// RateLimitAuthenticated limits requests per API key for authenticated endpoints.
func RateLimitAuthenticated(cache domain.CacheRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// API key hash is set by the Auth middleware
		keyHash, exists := c.Get("api_key_hash")
		if !exists {
			c.Next()
			return
		}

		rateLimit, _ := c.Get("api_key_rate_limit")
		limit, ok := rateLimit.(int)
		if !ok || limit <= 0 {
			limit = 60
		}

		key := fmt.Sprintf("rl:auth:%s", keyHash.(string))

		allowed, err := cache.RateLimit(c.Request.Context(), key, limit, rateLimitAuthWindow)
		if err != nil {
			logger.Log.Error("rate limit check failed", zap.Error(err))
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))

		if !allowed {
			errMsg := fmt.Sprintf("rate limit exceeded — %d requests/minute allowed", limit)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, domain.Response{
				Success: false,
				Data:    nil,
				Error:   &errMsg,
			})
			return
		}

		c.Next()
	}
}
