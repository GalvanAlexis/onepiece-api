package domain

import (
	"context"
	"time"
)

// ===== CHARACTER REPOSITORY =====

// CharacterFilters holds optional filters for character queries.
type CharacterFilters struct {
	CrewID *int64
	Status *CharacterStatus
	Limit  int
	Offset int
}

// CharacterRepository defines the interface for character data access.
type CharacterRepository interface {
	GetAll(ctx context.Context, filters CharacterFilters) ([]Character, int64, error)
	GetByID(ctx context.Context, id int64) (*Character, error)
}

// ===== ARC REPOSITORY =====

// ArcFilters holds optional filters for arc queries.
type ArcFilters struct {
	Saga   *string
	Limit  int
	Offset int
}

// ArcRepository defines the interface for arc data access.
type ArcRepository interface {
	GetAll(ctx context.Context, filters ArcFilters) ([]Arc, int64, error)
	GetByID(ctx context.Context, id int64) (*Arc, error)
}

// ===== DEVIL FRUIT REPOSITORY =====

// DevilFruitFilters holds optional filters for devil fruit queries.
type DevilFruitFilters struct {
	Type   *DevilFruitType
	Limit  int
	Offset int
}

// DevilFruitRepository defines the interface for devil fruit data access.
type DevilFruitRepository interface {
	GetAll(ctx context.Context, filters DevilFruitFilters) ([]DevilFruit, int64, error)
	GetByID(ctx context.Context, id int64) (*DevilFruit, error)
}

// ===== API KEY REPOSITORY =====

// APIKeyRepository defines the interface for API key data access.
type APIKeyRepository interface {
	Create(ctx context.Context, keyHash, label string, rateLimit int) (*APIKey, error)
	GetByHash(ctx context.Context, keyHash string) (*APIKey, error)
	UpdateLastUsed(ctx context.Context, id int64) error
}

// ===== CACHE REPOSITORY =====

// CacheRepository defines the interface for cache operations.
type CacheRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Del(ctx context.Context, key string) error
	// RateLimit implements sliding window rate limiting.
	// Returns (current count, error). Increments and returns false if limit exceeded.
	RateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
}
