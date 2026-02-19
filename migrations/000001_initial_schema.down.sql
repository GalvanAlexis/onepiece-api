-- ===========================
-- ONE PIECE API: Initial Schema
-- Migration: 000001 DOWN
-- ===========================

DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS character_haki;
DROP TABLE IF EXISTS chapters;
DROP TABLE IF EXISTS episodes;
DROP TABLE IF EXISTS characters;
DROP TABLE IF EXISTS haki_types;
DROP TABLE IF EXISTS devil_fruits;
DROP TABLE IF EXISTS arcs;
DROP TABLE IF EXISTS crews;

DROP EXTENSION IF EXISTS pg_trgm;
