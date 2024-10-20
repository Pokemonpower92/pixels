-- +migrate Up
CREATE TABLE authors (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    bio TEXT
);

-- +migrate Down
DROP TABLE authors;
