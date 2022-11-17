-- Filename: migrations/000003_add_toasts_indexes.down.sql
DROP INDEX If EXISTS toasts_name_idx;
DROP INDEX If EXISTS toasts_level_idx;
DROP INDEX If EXISTS toasts_mode_idx;