-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetByUserName :one
SELECT * FROM users
WHERE user_name = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at;

-- name: CreateUser :one
INSERT INTO users (
  id, user_name, permissions, created_at, updated_at
) VALUES (
  uuid_generate_v4(), $1, $2, NOW(), NOW() 
)
RETURNING *;

