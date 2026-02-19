package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/alexisgalvan/onepiece-api/internal/config"
	"github.com/alexisgalvan/onepiece-api/pkg/logger"
)

// NewRedisClient creates and validates a Redis client connection.
func NewRedisClient(cfg config.RedisConfig) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,

		// Connection pool settings
		PoolSize:        10,
		MinIdleConns:    2,
		ConnMaxLifetime: 1 * time.Hour,
		ConnMaxIdleTime: 30 * time.Minute,

		// Timeouts
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verify connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	logger.Log.Info("Redis connected",
		zap.String("addr", opts.Addr),
	)

	return client, nil
}
