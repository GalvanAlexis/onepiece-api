package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/alexisgalvan/onepiece-api/internal/config"
	"github.com/alexisgalvan/onepiece-api/internal/domain"
	"github.com/alexisgalvan/onepiece-api/internal/handler"
	"github.com/alexisgalvan/onepiece-api/internal/middleware"
)

// SetupRouter creates and configures the Gin engine.
func SetupRouter(
	cfg *config.Config,
	healthHandler *handler.HealthHandler,
	characterHandler *handler.CharacterHandler,
	apiKeyHandler *handler.APIKeyHandler,
	apiKeyRepo domain.APIKeyRepository,
	cacheRepo domain.CacheRepository,
) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))

	// ===== REGISTER ROUTES =====

	// Health — no auth, no rate limit
	router.GET("/health", healthHandler.Check)

	// /api/v1 group
	v1 := router.Group("/api/v1")
	{
		// Auth — generate API key (public, rate-limited by IP)
		auth := v1.Group("/auth")
		auth.Use(middleware.RateLimitPublic(cacheRepo, cfg.RateLimit.Public))
		{
			auth.POST("/api-key", apiKeyHandler.Generate)
		}

		// Public endpoints — rate limited by IP
		public := v1.Group("")
		public.Use(middleware.RateLimitPublic(cacheRepo, cfg.RateLimit.Public))
		{
			public.GET("/characters", characterHandler.List)
			public.GET("/characters/:id", characterHandler.GetByID)
		}

		// Authenticated endpoints — API key required + per-key rate limit (Phase 3+)
		protected := v1.Group("")
		protected.Use(middleware.Auth(apiKeyRepo, cacheRepo))
		protected.Use(middleware.RateLimitAuthenticated(cacheRepo))
		{
			_ = protected // Phase 3: write endpoints go here
		}
	}

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		errMsg := "route not found"
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   errMsg,
			"data":    nil,
		})
	})

	return router
}
