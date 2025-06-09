-- +migrate Up
CREATE TABLE images (
    db_id SERIAL PRIMARY KEY,
    id UUID UNIQUE NOT NULL,
    image_data BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);;
CREATE INDEX images_created_at ON images (created_at);
CREATE INDEX images_updated_at ON images (updated_at);
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +migrate Down
DROP TABLE images;
DROP EXTENSION "uuid-ossp";
