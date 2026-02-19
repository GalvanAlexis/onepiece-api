-- ===========================
-- ONE PIECE API: Initial Schema
-- Migration: 000001 UP
-- ===========================

-- Enable trigram extension for fuzzy search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- ===== CREWS =====
CREATE TABLE crews (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL UNIQUE,
    captain     VARCHAR(255) NOT NULL,
    affiliation VARCHAR(255),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ===== ARCS =====
CREATE TABLE arcs (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL UNIQUE,
    saga        VARCHAR(255) NOT NULL,
    description TEXT,
    order_index INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ===== DEVIL FRUITS =====
CREATE TABLE devil_fruits (
    id           BIGSERIAL PRIMARY KEY,
    name         VARCHAR(255) NOT NULL UNIQUE,
    type         VARCHAR(50) NOT NULL CHECK (type IN ('Logia', 'Zoan', 'Paramecia')),
    description  TEXT,
    current_user VARCHAR(255),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ===== HAKI TYPES =====
CREATE TABLE haki_types (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ===== CHARACTERS =====
CREATE TABLE characters (
    id             BIGSERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    alias          VARCHAR(255),
    status         VARCHAR(50) NOT NULL DEFAULT 'unknown' CHECK (status IN ('alive', 'deceased', 'unknown')),
    bounty         BIGINT,
    origin         VARCHAR(255),
    description    TEXT,
    image_url      TEXT,
    crew_id        BIGINT REFERENCES crews(id) ON DELETE SET NULL,
    devil_fruit_id BIGINT REFERENCES devil_fruits(id) ON DELETE SET NULL UNIQUE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ===== CHARACTER HAKI (Junction Table) =====
CREATE TABLE character_haki (
    character_id BIGINT NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    haki_id      BIGINT NOT NULL REFERENCES haki_types(id) ON DELETE CASCADE,
    PRIMARY KEY (character_id, haki_id)
);

-- ===== EPISODES =====
CREATE TABLE episodes (
    id         BIGSERIAL PRIMARY KEY,
    number     INT NOT NULL UNIQUE,
    title      VARCHAR(500) NOT NULL,
    arc_id     BIGINT REFERENCES arcs(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ===== CHAPTERS =====
CREATE TABLE chapters (
    id         BIGSERIAL PRIMARY KEY,
    number     INT NOT NULL UNIQUE,
    title      VARCHAR(500) NOT NULL,
    volume     INT,
    arc_id     BIGINT REFERENCES arcs(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ===== API KEYS =====
CREATE TABLE api_keys (
    id                BIGSERIAL PRIMARY KEY,
    key_hash          VARCHAR(64) NOT NULL UNIQUE,  -- SHA256 hex
    label             VARCHAR(255) NOT NULL,
    rate_limit_per_min INT NOT NULL DEFAULT 60,
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    last_used_at      TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ===========================
-- INDEXES
-- ===========================

-- Characters: full-text search on name
CREATE INDEX idx_characters_name_trgm ON characters USING GIN (name gin_trgm_ops);
CREATE INDEX idx_characters_crew_id ON characters (crew_id);
CREATE INDEX idx_characters_status ON characters (status);

-- Arcs: order and saga filtering
CREATE INDEX idx_arcs_saga ON arcs (saga);
CREATE INDEX idx_arcs_order ON arcs (order_index);

-- Devil fruits: type filtering
CREATE INDEX idx_devil_fruits_type ON devil_fruits (type);

-- Episodes: arc filtering
CREATE INDEX idx_episodes_arc_id ON episodes (arc_id);

-- Chapters: arc filtering
CREATE INDEX idx_chapters_arc_id ON chapters (arc_id);

-- API keys: fast lookup by hash
CREATE UNIQUE INDEX idx_api_keys_hash ON api_keys (key_hash);
CREATE INDEX idx_api_keys_active ON api_keys (is_active);

-- ===========================
-- SEED: Haki Types (static data)
-- ===========================
INSERT INTO haki_types (name, description) VALUES
    ('Kenbunshoku Haki', 'Observation Haki — allows the user to sense the presence, strength, and emotions of others.'),
    ('Busoshoku Haki', 'Armament Haki — allows the user to use their spirit as armor to attack or defend.'),
    ('Haoshoku Haki', 'Conqueror''s Haki — rare ability to exert one''s willpower over others.');
