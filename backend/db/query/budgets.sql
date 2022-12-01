-- name: CreateBudget :one
INSERT INTO budget (
budget_name, initial_amount, current_amount, description, start_date, end_date, user_id
) VALUES (
 $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetBudgetByID :one
SELECT * FROM budget 
WHERE id = $1 LIMIT 1;

-- name: GetBudgetByUserID :one
SELECT * FROM budget
WHERE user_id = $1 LIMIT 1;

-- name: ListBudgets :many
SELECT * FROM budget 
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateBudget :one
UPDATE budget
SET
budget_name = $2,
initial_amount = $3,
current_amount = $4,
description = $5,
start_date = $6,
end_date = $7,
user_id = $8,
updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteBudget :exec
DELETE FROM budget
WHERE id = $1;