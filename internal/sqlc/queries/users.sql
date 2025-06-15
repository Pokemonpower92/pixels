-- name: GetUser :one
SELECT * FROM users
WHERE user_name = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  id, user_name, password, created_at, updated_at
) VALUES (
  uuid_generate_v4(), $1, $2, NOW(), NOW() 
)
RETURNING *;

