-- +migrate Up
CREATE TABLE users (
    db_id SERIAL PRIMARY KEY,
    id UUID UNIQUE NOT NULL,
    user_name TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() 
);
CREATE INDEX users_id on users (id);
CREATE INDEX users_user_name ON users (user_name);

CREATE TABLE images (
    db_id SERIAL PRIMARY KEY,
    id UUID UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    image_data BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX images_id on images (id);
CREATE INDEX images_user_id ON images (user_id);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +migrate Down
DROP TABLE images;
DROP TABLE users;
DROP EXTENSION "uuid-ossp";
