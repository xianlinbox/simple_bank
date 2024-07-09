-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: AddSession :one
INSERT INTO sessions (
  id, username, refresh_token, user_agent, client_ip, expired_at
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;