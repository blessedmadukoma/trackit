-- name: CreateAccount :one
INSERT INTO accounts (
 user_id, balance
) VALUES (
 $1, $2
)
RETURNING *;

-- name: GetAccountByID :one
SELECT * FROM accounts 
WHERE id = $1 LIMIT 1;

-- name: GetAccountByUserID :one
SELECT * FROM accounts
WHERE user_id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts 
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET 
balance = $2,
updated_at = now()
WHERE id = $1
RETURNING *;
