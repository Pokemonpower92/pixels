-- name: GetImage :one
SELECT * FROM images
WHERE id = $1 LIMIT 1;

-- name: ListImages :many
SELECT * FROM images
ORDER BY created_at;

-- name: CreateImage :one
INSERT INTO images (
  id, image_data, created_at, updated_at
) VALUES (
  uuid_generate_v4(), $1, NOW(), NOW() 
)
RETURNING *;

