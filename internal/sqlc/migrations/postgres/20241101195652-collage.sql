-- +migrate Up
CREATE TABLE collages (
    db_id SERIAL PRIMARY KEY,
    id UUID UNIQUE NOT NULL,
    name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    image_set_id UUID NOT NULL,
    target_image_id UUID NOT NULL,
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL,
    FOREIGN KEY (image_set_id) REFERENCES image_sets(id),
    FOREIGN KEY (target_image_id) REFERENCES target_images(id)
);
CREATE INDEX collages_created_at ON collages (created_at);
CREATE INDEX collages_updated_at ON collages (updated_at);

CREATE TABLE collage_images (
    db_id SERIAL PRIMARY KEY,
    id UUID UNIQUE NOT NULL,
    collage_id UUID UNIQUE NOT NULL,
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL,
    FOREIGN KEY (collage_id) REFERENCES collages(id)
);
CREATE INDEX collage_images_created_at ON collage_images (created_at);
CREATE INDEX collage_images_updated_at ON collage_images (updated_at);

-- +migrate Down
DROP TABLE collages;
DROP TABLE collage_images;
