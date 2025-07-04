-- name: GetImage :one
SELECT * FROM images
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: ListImages :many
SELECT * FROM images
WHERE user_id = $1;

-- name: CreateImage :one
INSERT INTO images (
  id, user_id, image_data, created_at, updated_at
) VALUES (
  uuid_generate_v4(), $1, $2, NOW(), NOW() 
)
RETURNING *;

