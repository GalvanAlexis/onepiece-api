package redisrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheRepository implements domain.CacheRepository using Redis.
type CacheRepository struct {
	client *redis.Client
}

// NewCacheRepository creates a new Redis-backed CacheRepository.
func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

// Get retrieves a cached string value. Returns error if key not found.
func (r *CacheRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("cache miss: %s", key)
	}
	return val, err
}

// Set stores a string value with TTL.
func (r *CacheRepository) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Del removes a key from cache.
func (r *CacheRepository) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// RateLimit implements a Redis sliding window rate limiter using sorted sets.
// Returns true if the request is allowed, false if the limit has been exceeded.
func (r *CacheRepository) RateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	now := time.Now()
	windowStart := now.Add(-window).UnixMilli()
	nowMs := now.UnixMilli()
	member := fmt.Sprintf("%d", nowMs)

	pipe := r.client.Pipeline()

	// Remove expired entries outside the window
	pipe.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%d", windowStart))
	// Add current request
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(nowMs), Member: member})
	// Count requests in window
	pipe.ZCard(ctx, key)
	// Reset TTL
	pipe.Expire(ctx, key, window)

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("rate limit pipeline: %w", err)
	}

	// ZCard is the 3rd command (index 2)
	count := cmds[2].(*redis.IntCmd).Val()
	return count <= int64(limit), nil
}
