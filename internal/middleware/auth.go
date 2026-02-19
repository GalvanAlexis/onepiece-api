package middleware

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alexisgalvan/onepiece-api/internal/domain"
	"github.com/alexisgalvan/onepiece-api/pkg/logger"
)

const apiKeyCacheTTL = 5 * time.Minute

// Auth validates the X-API-KEY header and sets api_key_hash + api_key_rate_limit in context.
func Auth(apiKeyRepo domain.APIKeyRepository, cache domain.CacheRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawKey := c.GetHeader("X-API-KEY")
		if rawKey == "" {
			errMsg := "missing X-API-KEY header"
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.Response{
				Success: false,
				Data:    nil,
				Error:   &errMsg,
			})
			return
		}

		// Strip prefix if present
		rawKey = strings.TrimSpace(rawKey)

		// Hash the provided key
		hash := sha256.Sum256([]byte(rawKey))
		keyHash := fmt.Sprintf("%x", hash)

		// Try cache first
		var apiKey *domain.APIKey
		cacheKey := fmt.Sprintf("apikey:%s", keyHash)

		cached, err := cache.Get(c.Request.Context(), cacheKey)
		if err == nil && cached == "valid" {
			// Key is valid (cached), just set context
			c.Set("api_key_hash", keyHash)
			c.Next()
			return
		}

		// Cache miss — check DB
		apiKey, err = apiKeyRepo.GetByHash(c.Request.Context(), keyHash)
		if err != nil || apiKey == nil || !apiKey.IsActive {
			logger.Log.Warn("invalid API key attempt",
				zap.String("hash_prefix", keyHash[:8]),
				zap.String("ip", c.ClientIP()),
			)
			errMsg := "invalid or inactive API key"
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.Response{
				Success: false,
				Data:    nil,
				Error:   &errMsg,
			})
			return
		}

		// Cache valid key
		_ = cache.Set(c.Request.Context(), cacheKey, "valid", apiKeyCacheTTL)

		// Update last used async (don't block request)
		go func() {
			_ = apiKeyRepo.UpdateLastUsed(c.Request.Context(), apiKey.ID)
		}()

		// Set context values for rate limit middleware
		c.Set("api_key_hash", keyHash)
		c.Set("api_key_rate_limit", apiKey.RateLimitPerMin)
		c.Set("api_key_id", apiKey.ID)

		c.Next()
	}
}
