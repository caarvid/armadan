// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package schema

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email, password, role
`

type CreateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg *CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.Role,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, password, role FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.Role,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, email, password, role FROM users WHERE id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.Role,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, email, password, role FROM users
`

func (q *Queries) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.Password,
			&i.Role,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserEmail = `-- name: UpdateUserEmail :one
UPDATE users SET email = $1 WHERE id = $2 RETURNING id, email, password, role
`

type UpdateUserEmailParams struct {
	Email string    `json:"email"`
	ID    uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUserEmail(ctx context.Context, arg *UpdateUserEmailParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUserEmail, arg.Email, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.Role,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :one
UPDATE users SET password = $1 WHERE id = $2 RETURNING id, email, password, role
`

type UpdateUserPasswordParams struct {
	Password string    `json:"password"`
	ID       uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg *UpdateUserPasswordParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUserPassword, arg.Password, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.Role,
	)
	return i, err
}

const updateUserRole = `-- name: UpdateUserRole :one
UPDATE users SET role = $1 WHERE id = $2 RETURNING id, email, password, role
`

type UpdateUserRoleParams struct {
	Role UsersRoleEnum `json:"role"`
	ID   uuid.UUID     `json:"id"`
}

func (q *Queries) UpdateUserRole(ctx context.Context, arg *UpdateUserRoleParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUserRole, arg.Role, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.Role,
	)
	return i, err
}