// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: users.sql

package db

import (
	"context"
)

const createUserAccount = `-- name: CreateUserAccount :one
INSERT INTO users (
 firstname, lastname, email, mobile, password
) VALUES (
 $1, $2, $3, $4, $5
)
RETURNING id, firstname, lastname, email, mobile, password, created_at, updated_at
`

type CreateUserAccountParams struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
	Password  string `json:"password"`
}

func (q *Queries) CreateUserAccount(ctx context.Context, arg CreateUserAccountParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUserAccount,
		arg.Firstname,
		arg.Lastname,
		arg.Email,
		arg.Mobile,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Firstname,
		&i.Lastname,
		&i.Email,
		&i.Mobile,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUserAccount = `-- name: DeleteUserAccount :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUserAccount(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUserAccount, id)
	return err
}

const getUserAccountByID = `-- name: GetUserAccountByID :one
SELECT id, firstname, lastname, email, mobile, password, created_at, updated_at FROM users 
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserAccountByID(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserAccountByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Firstname,
		&i.Lastname,
		&i.Email,
		&i.Mobile,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, firstname, lastname, email, mobile, password, created_at, updated_at FROM users 
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Firstname,
			&i.Lastname,
			&i.Email,
			&i.Mobile,
			&i.Password,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateUserAccount = `-- name: UpdateUserAccount :one
UPDATE users
SET 
firstname = $2,
lastname = $3,
email = $4,
mobile = $5,
password = $6,
updated_at = now()
WHERE id = $1
RETURNING id, firstname, lastname, email, mobile, password, created_at, updated_at
`

type UpdateUserAccountParams struct {
	ID        int64  `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
	Password  string `json:"password"`
}

func (q *Queries) UpdateUserAccount(ctx context.Context, arg UpdateUserAccountParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserAccount,
		arg.ID,
		arg.Firstname,
		arg.Lastname,
		arg.Email,
		arg.Mobile,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Firstname,
		&i.Lastname,
		&i.Email,
		&i.Mobile,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}