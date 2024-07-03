-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;


-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY id;

-- name: AddTransfer :one
INSERT INTO transfers (
  from_account, to_account, amount
) VALUES (
  $1, $2, $3
)
RETURNING *;
