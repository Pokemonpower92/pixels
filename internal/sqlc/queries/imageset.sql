-- name: GetImageset :one
SELECT * FROM image_sets
WHERE id = $1 LIMIT 1;

-- name: Listimage_sets :many
SELECT * FROM image_sets
ORDER BY name;

-- name: CreateImageset :one
INSERT INTO image_sets (
  id, name, description, created_at, updated_at
) VALUES (
  uuid_generate_v4(), $1, $2, NOW(), NOW() 
)
RETURNING *;

