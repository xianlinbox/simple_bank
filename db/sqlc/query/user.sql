-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: AddUser :one
INSERT INTO users (
  username, password, full_name, email, role
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET 
  password = coalesce(sqlc.narg(password), password), 
  full_name = coalesce(sqlc.narg(full_name), full_name), 
  email = coalesce(sqlc.narg(email), email),
  password_expired_at = coalesce(sqlc.narg(password_expired_at), password_expired_at),
  role = coalesce(sqlc.narg(role), role)
WHERE username = sqlc.arg(username)
RETURNING *;