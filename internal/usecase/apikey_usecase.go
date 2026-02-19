package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/alexisgalvan/onepiece-api/internal/domain"
)

// APIKeyUseCase handles the business logic for API key management.
type APIKeyUseCase struct {
	repo domain.APIKeyRepository
}

// NewAPIKeyUseCase creates a new APIKeyUseCase.
func NewAPIKeyUseCase(repo domain.APIKeyRepository) *APIKeyUseCase {
	return &APIKeyUseCase{repo: repo}
}

// CreateResult holds the newly generated API key (only exposed once).
type CreateResult struct {
	Key  string        `json:"key"` // plain text — shown only once
	Meta domain.APIKey `json:"meta"`
}

// Generate creates a new API key with the given label and rate limit.
// The plain-text key is returned only once; only the hash is stored.
func (uc *APIKeyUseCase) Generate(ctx context.Context, label string, rateLimit int) (*CreateResult, error) {
	if label == "" {
		label = "default"
	}
	if rateLimit <= 0 {
		rateLimit = 60
	}

	// Generate 32 cryptographically random bytes
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("generate random bytes: %w", err)
	}

	// Encode to base64url and add prefix
	rawKey := "op_" + base64.RawURLEncoding.EncodeToString(b)

	// Hash with SHA256 (only stored value)
	h := sha256.Sum256([]byte(rawKey))
	keyHash := fmt.Sprintf("%x", h)

	// Persist to DB
	apiKey, err := uc.repo.Create(ctx, keyHash, label, rateLimit)
	if err != nil {
		return nil, fmt.Errorf("persist api key: %w", err)
	}

	return &CreateResult{
		Key:  rawKey,
		Meta: *apiKey,
	}, nil
}
