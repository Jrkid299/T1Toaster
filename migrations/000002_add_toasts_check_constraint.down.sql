-- Filename: migrations/000002_add_toasts_check_constraint.down.sql

ALTER TABLE toasts DROP CONSTRAINT IF EXISTS mode_length_check;