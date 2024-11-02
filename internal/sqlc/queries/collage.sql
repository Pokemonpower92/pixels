-- name: GetCollage :one
SELECT * FROM collages
WHERE id = $1 LIMIT 1;

-- name: ListCollages :many
SELECT * FROM collages
ORDER BY name;

-- name: CreateCollage :one
INSERT INTO collages (
  id, name, description, image_set_id, target_image_id, created_at, updated_at
) VALUES (
  uuid_generate_v4(), $1, $2, $3, $4, NOW(), NOW() 
)
RETURNING *;

