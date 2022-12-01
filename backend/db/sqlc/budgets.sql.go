// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: budgets.sql

package db

import (
	"context"
	"time"
)

const createBudget = `-- name: CreateBudget :one
INSERT INTO budget (
budget_name, initial_amount, current_amount, description, start_date, end_date, user_id
) VALUES (
 $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, created_at, updated_at, budget_name, initial_amount, current_amount, description, start_date, end_date, user_id
`

type CreateBudgetParams struct {
	BudgetName    string    `json:"budget_name"`
	InitialAmount float64   `json:"initial_amount"`
	CurrentAmount float64   `json:"current_amount"`
	Description   string    `json:"description"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	UserID        int64     `json:"user_id"`
}

func (q *Queries) CreateBudget(ctx context.Context, arg CreateBudgetParams) (Budget, error) {
	row := q.db.QueryRowContext(ctx, createBudget,
		arg.BudgetName,
		arg.InitialAmount,
		arg.CurrentAmount,
		arg.Description,
		arg.StartDate,
		arg.EndDate,
		arg.UserID,
	)
	var i Budget
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.BudgetName,
		&i.InitialAmount,
		&i.CurrentAmount,
		&i.Description,
		&i.StartDate,
		&i.EndDate,
		&i.UserID,
	)
	return i, err
}

const deleteBudget = `-- name: DeleteBudget :exec
DELETE FROM budget
WHERE id = $1
`

func (q *Queries) DeleteBudget(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteBudget, id)
	return err
}

const getBudgetByID = `-- name: GetBudgetByID :one
SELECT id, created_at, updated_at, budget_name, initial_amount, current_amount, description, start_date, end_date, user_id FROM budget 
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetBudgetByID(ctx context.Context, id int64) (Budget, error) {
	row := q.db.QueryRowContext(ctx, getBudgetByID, id)
	var i Budget
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.BudgetName,
		&i.InitialAmount,
		&i.CurrentAmount,
		&i.Description,
		&i.StartDate,
		&i.EndDate,
		&i.UserID,
	)
	return i, err
}

const getBudgetByUserID = `-- name: GetBudgetByUserID :one
SELECT id, created_at, updated_at, budget_name, initial_amount, current_amount, description, start_date, end_date, user_id FROM budget
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetBudgetByUserID(ctx context.Context, userID int64) (Budget, error) {
	row := q.db.QueryRowContext(ctx, getBudgetByUserID, userID)
	var i Budget
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.BudgetName,
		&i.InitialAmount,
		&i.CurrentAmount,
		&i.Description,
		&i.StartDate,
		&i.EndDate,
		&i.UserID,
	)
	return i, err
}

const listBudgets = `-- name: ListBudgets :many
SELECT id, created_at, updated_at, budget_name, initial_amount, current_amount, description, start_date, end_date, user_id FROM budget 
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListBudgetsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListBudgets(ctx context.Context, arg ListBudgetsParams) ([]Budget, error) {
	rows, err := q.db.QueryContext(ctx, listBudgets, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Budget
	for rows.Next() {
		var i Budget
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.BudgetName,
			&i.InitialAmount,
			&i.CurrentAmount,
			&i.Description,
			&i.StartDate,
			&i.EndDate,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBudget = `-- name: UpdateBudget :one
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
RETURNING id, created_at, updated_at, budget_name, initial_amount, current_amount, description, start_date, end_date, user_id
`

type UpdateBudgetParams struct {
	ID            int64     `json:"id"`
	BudgetName    string    `json:"budget_name"`
	InitialAmount float64   `json:"initial_amount"`
	CurrentAmount float64   `json:"current_amount"`
	Description   string    `json:"description"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	UserID        int64     `json:"user_id"`
}

func (q *Queries) UpdateBudget(ctx context.Context, arg UpdateBudgetParams) (Budget, error) {
	row := q.db.QueryRowContext(ctx, updateBudget,
		arg.ID,
		arg.BudgetName,
		arg.InitialAmount,
		arg.CurrentAmount,
		arg.Description,
		arg.StartDate,
		arg.EndDate,
		arg.UserID,
	)
	var i Budget
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.BudgetName,
		&i.InitialAmount,
		&i.CurrentAmount,
		&i.Description,
		&i.StartDate,
		&i.EndDate,
		&i.UserID,
	)
	return i, err
}
