-- name: GetCollageImage :one
SELECT * FROM collage_images
WHERE id = $1 LIMIT 1;

-- name: GetByCollageId :one
SELECT * FROM collage_images
WHERE collage_id = $1 LIMIT 1;

-- name: ListCollageImages :many
SELECT * FROM collage_images
ORDER BY updated_at;

-- name: CreateCollageImage :one
INSERT INTO collage_images (
  id, collage_id, created_at, updated_at
) VALUES (
  uuid_generate_v4(), $1, NOW(), NOW() 
)
RETURNING *;

