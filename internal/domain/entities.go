package domain

import "time"

// ===== CREW =====

// Crew represents a pirate crew in the One Piece universe.
type Crew struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Captain     string    `json:"captain" db:"captain"`
	Affiliation string    `json:"affiliation,omitempty" db:"affiliation"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ===== ARC =====

// Arc represents a story arc in the One Piece series.
type Arc struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Saga        string    `json:"saga" db:"saga"`
	Description string    `json:"description,omitempty" db:"description"`
	OrderIndex  int       `json:"order_index" db:"order_index"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ===== DEVIL FRUIT =====

// DevilFruitType represents the classification of a devil fruit.
type DevilFruitType string

const (
	DevilFruitLogia   DevilFruitType = "Logia"
	DevilFruitZoan    DevilFruitType = "Zoan"
	DevilFruitParamecia DevilFruitType = "Paramecia"
)

// DevilFruit represents a devil fruit in the One Piece universe.
type DevilFruit struct {
	ID          int64          `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Type        DevilFruitType `json:"type" db:"type"`
	Description string         `json:"description,omitempty" db:"description"`
	CurrentUser *string        `json:"current_user,omitempty" db:"current_user"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

// ===== HAKI =====

// HakiType represents a type of Haki.
type HakiType struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// ===== CHARACTER =====

// CharacterStatus represents a character's current status.
type CharacterStatus string

const (
	CharacterAlive   CharacterStatus = "alive"
	CharacterDeceased CharacterStatus = "deceased"
	CharacterUnknown CharacterStatus = "unknown"
)

// Character represents a character in the One Piece universe.
type Character struct {
	ID           int64           `json:"id" db:"id"`
	Name         string          `json:"name" db:"name"`
	Alias        *string         `json:"alias,omitempty" db:"alias"`
	Status       CharacterStatus `json:"status" db:"status"`
	Bounty       *int64          `json:"bounty,omitempty" db:"bounty"`
	Origin       *string         `json:"origin,omitempty" db:"origin"`
	Description  *string         `json:"description,omitempty" db:"description"`
	ImageURL     *string         `json:"image_url,omitempty" db:"image_url"`
	CrewID       *int64          `json:"crew_id,omitempty" db:"crew_id"`
	DevilFruitID *int64          `json:"devil_fruit_id,omitempty" db:"devil_fruit_id"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`

	// Joined fields (populated in queries)
	Crew       *Crew       `json:"crew,omitempty"`
	DevilFruit *DevilFruit `json:"devil_fruit,omitempty"`
	HakiTypes  []HakiType  `json:"haki_types,omitempty"`
}

// ===== EPISODE =====

// Episode represents an anime episode.
type Episode struct {
	ID        int64     `json:"id" db:"id"`
	Number    int       `json:"number" db:"number"`
	Title     string    `json:"title" db:"title"`
	ArcID     *int64    `json:"arc_id,omitempty" db:"arc_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ===== CHAPTER =====

// Chapter represents a manga chapter.
type Chapter struct {
	ID        int64     `json:"id" db:"id"`
	Number    int       `json:"number" db:"number"`
	Title     string    `json:"title" db:"title"`
	Volume    *int      `json:"volume,omitempty" db:"volume"`
	ArcID     *int64    `json:"arc_id,omitempty" db:"arc_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ===== API KEY =====

// APIKey represents an API authentication key.
type APIKey struct {
	ID             int64     `json:"id" db:"id"`
	KeyHash        string    `json:"-" db:"key_hash"` // Never exposed in JSON
	Label          string    `json:"label" db:"label"`
	RateLimitPerMin int      `json:"rate_limit_per_min" db:"rate_limit_per_min"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	LastUsedAt     *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// ===== RESPONSE STRUCTS =====

// PaginationMeta holds pagination metadata for list responses.
type PaginationMeta struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasNext bool `json:"has_next"`
}

// Response is the standard API response envelope.
type Response struct {
	Success bool            `json:"success"`
	Data    interface{}     `json:"data"`
	Error   *string         `json:"error"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
}
