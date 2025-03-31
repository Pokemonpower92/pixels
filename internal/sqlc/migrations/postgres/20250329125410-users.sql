-- Adds the users table to the database.

-- +migrate Up
CREATE TABLE users (
    db_id SERIAL PRIMARY KEY,
    id UUID UNIQUE NOT NULL,
    user_name TEXT UNIQUE NOT NULL,
    permissions JSONB DEFAULT '[]'::jsonb,
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL
);
CREATE INDEX users_user_name ON users (user_name);
CREATE INDEX users_created_at ON users (created_at);
CREATE INDEX users_updated_at ON users (updated_at);

-- +migrate Down
DROP TABLE users;