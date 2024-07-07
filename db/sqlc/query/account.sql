-- name: GetAccount :one
SELECT * FROM account
WHERE id = $1 LIMIT 1;

-- name: GetAccountsByOwner :many
SELECT * FROM account
WHERE owner = $1
ORDER BY id
LIMIT $2;

-- name: AddAccount :one
INSERT INTO account (
  owner, balance, currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateAccount :exec
UPDATE account
  set balance = $2
WHERE id = $1;

-- name: DeleteAccount :exec
DELETE FROM account
WHERE id = $1;
