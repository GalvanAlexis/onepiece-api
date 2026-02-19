package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alexisgalvan/onepiece-api/internal/config"
	"github.com/alexisgalvan/onepiece-api/internal/handler"
	"github.com/alexisgalvan/onepiece-api/internal/middleware"
	postgresrepo "github.com/alexisgalvan/onepiece-api/internal/repository/postgres"
	redisrepo "github.com/alexisgalvan/onepiece-api/internal/repository/redis"
	"github.com/alexisgalvan/onepiece-api/internal/usecase"
	"github.com/alexisgalvan/onepiece-api/pkg/database"
	"github.com/alexisgalvan/onepiece-api/pkg/logger"
)

// @title           One Piece API
// @version         1.0.0
// @description     A public REST API for the One Piece universe. Inspired by PokéAPI.
// @contact.name    Alexis Galvan
// @license.name    MIT
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	// ===== 1. LOAD CONFIG =====
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// ===== 2. INIT LOGGER =====
	if err := logger.Init(cfg.App.Env, cfg.Log.Level); err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Log.Info("Starting One Piece API",
		zap.String("version", cfg.App.Version),
		zap.String("env", cfg.App.Env),
		zap.String("port", cfg.App.Port),
	)

	// ===== 3. CONNECT POSTGRESQL =====
	pgPool, err := database.NewPostgresPool(cfg.Database)
	if err != nil {
		logger.Log.Fatal("failed to connect to PostgreSQL", zap.Error(err))
	}
	defer pgPool.Close()

	// ===== 4. CONNECT REDIS =====
	redisClient, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		logger.Log.Fatal("failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()

	// ===== 5. WIRE DEPENDENCIES =====
	// Repositories
	characterRepo := postgresrepo.NewCharacterRepository(pgPool)
	apiKeyRepo := postgresrepo.NewAPIKeyRepository(pgPool)
	cacheRepo := redisrepo.NewCacheRepository(redisClient)

	// Use cases
	characterUC := usecase.NewCharacterUseCase(characterRepo, cacheRepo)
	apiKeyUC := usecase.NewAPIKeyUseCase(apiKeyRepo)

	// Handlers
	healthHandler := handler.NewHealthHandler(cfg)
	characterHandler := handler.NewCharacterHandler(characterUC)
	apiKeyHandler := handler.NewAPIKeyHandler(apiKeyUC)

	// ===== 6. SETUP GIN =====
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

	// ===== 7. REGISTER ROUTES =====

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

	// ===== 8. START SERVER =====
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		logger.Log.Info("Server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("server error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Log.Info("Server stopped gracefully")
}
