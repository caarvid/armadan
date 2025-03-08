// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package schema

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, email, password) VALUES (?, ?, ?) RETURNING id, email, password, user_role
`

type CreateUserParams struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg *CreateUserParams) (User, error) {
	row := q.queryRow(ctx, q.createUserStmt, createUser, arg.ID, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.UserRole,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, id string) error {
	_, err := q.exec(ctx, q.deleteUserStmt, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, password, user_role FROM users WHERE email = ?
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.queryRow(ctx, q.getUserByEmailStmt, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.UserRole,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, email, password, user_role FROM users WHERE id = ?
`

func (q *Queries) GetUserById(ctx context.Context, id string) (User, error) {
	row := q.queryRow(ctx, q.getUserByIdStmt, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.UserRole,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, email, password, user_role FROM users
`

func (q *Queries) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := q.query(ctx, q.getUsersStmt, getUsers)
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
			&i.UserRole,
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

const updateUserEmail = `-- name: UpdateUserEmail :one
UPDATE users SET email = ? WHERE id = ? RETURNING id, email, password, user_role
`

type UpdateUserEmailParams struct {
	Email string `json:"email"`
	ID    string `json:"id"`
}

func (q *Queries) UpdateUserEmail(ctx context.Context, arg *UpdateUserEmailParams) (User, error) {
	row := q.queryRow(ctx, q.updateUserEmailStmt, updateUserEmail, arg.Email, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.UserRole,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :one
UPDATE users SET password = ? WHERE id = ? RETURNING id, email, password, user_role
`

type UpdateUserPasswordParams struct {
	Password string `json:"password"`
	ID       string `json:"id"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg *UpdateUserPasswordParams) (User, error) {
	row := q.queryRow(ctx, q.updateUserPasswordStmt, updateUserPassword, arg.Password, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.UserRole,
	)
	return i, err
}

const updateUserRole = `-- name: UpdateUserRole :one
UPDATE users SET user_role = ? WHERE id = ? RETURNING id, email, password, user_role
`

type UpdateUserRoleParams struct {
	UserRole string `json:"userRole"`
	ID       string `json:"id"`
}

func (q *Queries) UpdateUserRole(ctx context.Context, arg *UpdateUserRoleParams) (User, error) {
	row := q.queryRow(ctx, q.updateUserRoleStmt, updateUserRole, arg.UserRole, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.UserRole,
	)
	return i, err
}
