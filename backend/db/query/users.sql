-- name: CreateUserAccount :one
INSERT INTO users (
 firstname, lastname, email, mobile, password
) VALUES (
 $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUserAccountByID :one
SELECT * FROM users 
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users 
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUserAccount :one
UPDATE users
SET 
firstname = $2,
lastname = $3,
email = $4,
mobile = $5,
password = $6,
updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUserAccount :exec
DELETE FROM users
WHERE id = $1;