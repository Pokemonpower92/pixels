
-- +migrate Up
CREATE TABLE target_images (
    db_id SERIAL PRIMARY KEY,
    id UUID UNIQUE NOT NULL,
    name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL
);
CREATE INDEX target_images_created_at ON target_images (created_at);
CREATE INDEX target_images_updated_at ON target_images (updated_at);

-- +migrate Down
DROP TABLE target_images;
