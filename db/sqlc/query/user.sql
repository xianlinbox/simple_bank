-- name: GetUser :one
SELECT * FROM user
WHERE username = $1 LIMIT 1;

-- -- name: AddUser :one
-- INSERT INTO user (
--   username, password, full_name, email
-- ) VALUES (
--   $1, $2, $3, $4
-- )
-- RETURNING *;