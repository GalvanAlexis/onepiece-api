package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/alexisgalvan/onepiece-api/internal/domain"
)

// APIKeyRepository implements domain.APIKeyRepository using PostgreSQL.
type APIKeyRepository struct {
	db *pgxpool.Pool
}

// NewAPIKeyRepository creates a new APIKeyRepository.
func NewAPIKeyRepository(db *pgxpool.Pool) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

// Create inserts a new API key record.
func (r *APIKeyRepository) Create(ctx context.Context, keyHash, label string, rateLimit int) (*domain.APIKey, error) {
	query := `
		INSERT INTO api_keys (key_hash, label, rate_limit_per_min)
		VALUES ($1, $2, $3)
		RETURNING id, key_hash, label, rate_limit_per_min, is_active, last_used_at, created_at
	`

	var key domain.APIKey
	row := r.db.QueryRow(ctx, query, keyHash, label, rateLimit)
	err := row.Scan(
		&key.ID,
		&key.KeyHash,
		&key.Label,
		&key.RateLimitPerMin,
		&key.IsActive,
		&key.LastUsedAt,
		&key.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create api key: %w", err)
	}

	return &key, nil
}

// GetByHash retrieves an active API key by its SHA256 hash.
func (r *APIKeyRepository) GetByHash(ctx context.Context, keyHash string) (*domain.APIKey, error) {
	query := `
		SELECT id, key_hash, label, rate_limit_per_min, is_active, last_used_at, created_at
		FROM api_keys
		WHERE key_hash = $1 AND is_active = true
	`

	var key domain.APIKey
	row := r.db.QueryRow(ctx, query, keyHash)
	err := row.Scan(
		&key.ID,
		&key.KeyHash,
		&key.Label,
		&key.RateLimitPerMin,
		&key.IsActive,
		&key.LastUsedAt,
		&key.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get api key by hash: %w", err)
	}

	return &key, nil
}

// UpdateLastUsed sets last_used_at to the current timestamp.
func (r *APIKeyRepository) UpdateLastUsed(ctx context.Context, id int64) error {
	query := `UPDATE api_keys SET last_used_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("update api key last used: %w", err)
	}
	return nil
}
