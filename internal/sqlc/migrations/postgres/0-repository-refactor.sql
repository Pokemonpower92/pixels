-- +migrate Up
CREATE TABLE image_sets (
    db_id SERIAL PRIMARY KEY,
    id UUID UNIQUE NOT NULL,
    name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL
);
CREATE INDEX image_sets_created_at ON image_sets (created_at);
CREATE INDEX image_sets_updated_at ON image_sets (updated_at);

CREATE TABLE average_colors (
    db_id SERIAL PRIMARY KEY,
    id UUID UNIQUE NOT NULL,
    imageset_id UUID NOT NULL,
    file_name TEXT NOT NULL,
    r INTEGER NOT NULL,
    g INTEGER NOT NULL,
    b INTEGER NOT NULL,
    a INTEGER NOT NULL,
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL,
    FOREIGN KEY (imageset_id) REFERENCES image_sets(id)
);
CREATE INDEX average_colors_created_at ON average_colors (created_at);
CREATE INDEX average_colors_updated_at ON average_colors (updated_at);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +migrate Down
DROP TABLE image_sets;
DROP TABLE average_colors;
DROP EXTENSION "uuid-oosp";
