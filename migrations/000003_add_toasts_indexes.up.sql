-- Filename: migrations/000003_add_toasts_indexes.up.sql
CREATE INDEX IF NOT EXISTS toasts_name_idx ON toasts USING GIN(to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS toasts_level_idx ON toasts USING GIN(to_tsvector('simple', level));
CREATE INDEX IF NOT EXISTS toasts_mode_idx ON toasts USING GIN(mode);