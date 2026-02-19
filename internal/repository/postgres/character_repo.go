package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/alexisgalvan/onepiece-api/internal/domain"
)

// CharacterRepository implements domain.CharacterRepository using PostgreSQL.
type CharacterRepository struct {
	db *pgxpool.Pool
}

// NewCharacterRepository creates a new CharacterRepository.
func NewCharacterRepository(db *pgxpool.Pool) *CharacterRepository {
	return &CharacterRepository{db: db}
}

// GetAll retrieves a paginated list of characters with optional filters.
func (r *CharacterRepository) GetAll(ctx context.Context, filters domain.CharacterFilters) ([]domain.Character, int64, error) {
	args := []interface{}{}
	argIndex := 1
	where := "WHERE 1=1"

	if filters.CrewID != nil {
		where += fmt.Sprintf(" AND c.crew_id = $%d", argIndex)
		args = append(args, *filters.CrewID)
		argIndex++
	}

	if filters.Status != nil {
		where += fmt.Sprintf(" AND c.status = $%d", argIndex)
		args = append(args, string(*filters.Status))
		argIndex++
	}

	// Count total matching rows
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM characters c %s", where)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("character count: %w", err)
	}

	// Main query with pagination
	limit := filters.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := filters.Offset
	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf(`
		SELECT
			c.id, c.name, c.alias, c.status, c.bounty,
			c.origin, c.description, c.image_url,
			c.crew_id, c.devil_fruit_id,
			c.created_at, c.updated_at,
			cr.id, cr.name, cr.captain, cr.affiliation,
			df.id, df.name, df.type
		FROM characters c
		LEFT JOIN crews cr ON cr.id = c.crew_id
		LEFT JOIN devil_fruits df ON df.id = c.devil_fruit_id
		%s
		ORDER BY c.id ASC
		LIMIT $%d OFFSET $%d
	`, where, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("character list query: %w", err)
	}
	defer rows.Close()

	characters, err := scanCharacterRows(rows)
	if err != nil {
		return nil, 0, err
	}

	return characters, total, nil
}

// GetByID retrieves a single character by ID including Crew and DevilFruit.
func (r *CharacterRepository) GetByID(ctx context.Context, id int64) (*domain.Character, error) {
	query := `
		SELECT
			c.id, c.name, c.alias, c.status, c.bounty,
			c.origin, c.description, c.image_url,
			c.crew_id, c.devil_fruit_id,
			c.created_at, c.updated_at,
			cr.id, cr.name, cr.captain, cr.affiliation,
			df.id, df.name, df.type
		FROM characters c
		LEFT JOIN crews cr ON cr.id = c.crew_id
		LEFT JOIN devil_fruits df ON df.id = c.devil_fruit_id
		WHERE c.id = $1
	`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("character by id query: %w", err)
	}
	defer rows.Close()

	results, err := scanCharacterRows(rows)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return &results[0], nil
}

// scanCharacterRows scans character rows including optional crew/devil_fruit joins.
func scanCharacterRows(rows pgx.Rows) ([]domain.Character, error) {
	var characters []domain.Character

	for rows.Next() {
		var c domain.Character
		var crewID *int64
		var crewName, crewCaptain, crewAffil *string
		var dfID *int64
		var dfName, dfType *string

		err := rows.Scan(
			&c.ID, &c.Name, &c.Alias, &c.Status, &c.Bounty,
			&c.Origin, &c.Description, &c.ImageURL,
			&c.CrewID, &c.DevilFruitID,
			&c.CreatedAt, &c.UpdatedAt,
			&crewID, &crewName, &crewCaptain, &crewAffil,
			&dfID, &dfName, &dfType,
		)
		if err != nil {
			return nil, fmt.Errorf("scan character row: %w", err)
		}

		if crewID != nil {
			c.Crew = &domain.Crew{
				ID:      *crewID,
				Name:    *crewName,
				Captain: *crewCaptain,
				Affiliation: func() string {
					if crewAffil != nil {
						return *crewAffil
					}
					return ""
				}(),
			}
		}

		if dfID != nil {
			c.DevilFruit = &domain.DevilFruit{
				ID:   *dfID,
				Name: *dfName,
				Type: domain.DevilFruitType(*dfType),
			}
		}

		characters = append(characters, c)
	}

	return characters, rows.Err()
}
