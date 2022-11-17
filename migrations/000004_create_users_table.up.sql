-- Filename: migrations/000004_create_toasts_users.up.sq1

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL, 
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    website text NOT NULL,
    activated bool NOT NUll,
    version integer NOT NULL DEFAULT 1
);