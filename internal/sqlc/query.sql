-- name: GetImageset :one
SELECT * FROM imagesets
WHERE id = $1 LIMIT 1;

-- name: ListImageset :many
SELECT * FROM imagesets
ORDER BY name;

-- name: CreateImageset :one
INSERT INTO imagesets (
  id, name, description, created_at, updated_at
) VALUES (
  uuid_generate_v4(), $1, $2, NOW(), NOW() 
)
RETURNING *;

