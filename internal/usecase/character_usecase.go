package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alexisgalvan/onepiece-api/internal/domain"
)

const (
	characterListCacheTTL   = 5 * time.Minute
	characterDetailCacheTTL = 10 * time.Minute
)

// CharacterUseCase handles the business logic for characters.
type CharacterUseCase struct {
	repo  domain.CharacterRepository
	cache domain.CacheRepository
}

// NewCharacterUseCase creates a new CharacterUseCase.
func NewCharacterUseCase(repo domain.CharacterRepository, cache domain.CacheRepository) *CharacterUseCase {
	return &CharacterUseCase{repo: repo, cache: cache}
}

// ListResult holds the result of a list operation.
type ListResult struct {
	Items []domain.Character
	Total int64
}

// List retrieves a paginated, filtered list of characters. Results are cached.
func (uc *CharacterUseCase) List(ctx context.Context, filters domain.CharacterFilters) (*ListResult, error) {
	cacheKey := buildCharacterListKey(filters)

	// Try cache
	if cached, err := uc.cache.Get(ctx, cacheKey); err == nil {
		var result ListResult
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return &result, nil
		}
	}

	// Fetch from DB
	characters, total, err := uc.repo.GetAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("list characters: %w", err)
	}

	result := &ListResult{Items: characters, Total: total}

	// Cache the result
	if data, err := json.Marshal(result); err == nil {
		_ = uc.cache.Set(ctx, cacheKey, string(data), characterListCacheTTL)
	}

	return result, nil
}

// GetByID retrieves a single character. Results are cached.
func (uc *CharacterUseCase) GetByID(ctx context.Context, id int64) (*domain.Character, error) {
	cacheKey := fmt.Sprintf("character:%d", id)

	// Try cache
	if cached, err := uc.cache.Get(ctx, cacheKey); err == nil {
		var character domain.Character
		if err := json.Unmarshal([]byte(cached), &character); err == nil {
			return &character, nil
		}
	}

	// Fetch from DB
	character, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get character: %w", err)
	}

	if character == nil {
		return nil, nil
	}

	// Cache the result
	if data, err := json.Marshal(character); err == nil {
		_ = uc.cache.Set(ctx, cacheKey, string(data), characterDetailCacheTTL)
	}

	return character, nil
}

func buildCharacterListKey(f domain.CharacterFilters) string {
	crewID := 0
	if f.CrewID != nil {
		crewID = int(*f.CrewID)
	}
	status := ""
	if f.Status != nil {
		status = string(*f.Status)
	}
	return fmt.Sprintf("characters:list:crew=%d:status=%s:limit=%d:offset=%d",
		crewID, status, f.Limit, f.Offset)
}
