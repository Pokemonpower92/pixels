-- +migrate Up
CREATE TABLE imagesets (
    db_id SERIAL PRIMARY KEY,
    id UUID UNIQUE NOT NULL,
    name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL
);
CREATE INDEX imagesets_created_at ON imagesets (created_at);
CREATE INDEX imagesets_updated_at ON imagesets (updated_at);

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
    FOREIGN KEY (imageset_id) REFERENCES imagesets(id)
);
CREATE INDEX average_colors_created_at ON average_colors (created_at);
CREATE INDEX average_colors_updated_at ON average_colors (updated_at);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +migrate Down
DROP TABLE imagesets;
DROP TABLE average_colors;
DROP EXTENSION "uuid-oosp";
