// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: players.sql

package schema

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const createPlayer = `-- name: CreatePlayer :one
INSERT INTO players (first_name, last_name, hcp, user_id) VALUES ($1, $2, $3, $4) RETURNING id, first_name, last_name, points, user_id, hcp
`

type CreatePlayerParams struct {
	FirstName string          `json:"firstName"`
	LastName  string          `json:"lastName"`
	Hcp       decimal.Decimal `json:"hcp"`
	UserID    uuid.UUID       `json:"userId"`
}

func (q *Queries) CreatePlayer(ctx context.Context, arg *CreatePlayerParams) (Player, error) {
	row := q.db.QueryRow(ctx, createPlayer,
		arg.FirstName,
		arg.LastName,
		arg.Hcp,
		arg.UserID,
	)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Points,
		&i.UserID,
		&i.Hcp,
	)
	return i, err
}

const deletePlayer = `-- name: DeletePlayer :exec
DELETE FROM players WHERE id = $1
`

func (q *Queries) DeletePlayer(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deletePlayer, id)
	return err
}

const getLeaderboard = `-- name: GetLeaderboard :many
SELECT
  p.id,
  p.first_name,
  p.last_name,
  p.points,
  COUNT(rounds.id) as nr_of_rounds
FROM players p
LEFT JOIN rounds ON rounds.player_id = p.id
GROUP BY p.id
ORDER BY p.points DESC, nr_of_rounds DESC
`

type GetLeaderboardRow struct {
	ID         uuid.UUID `json:"id"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Points     int32     `json:"points"`
	NrOfRounds int64     `json:"nrOfRounds"`
}

func (q *Queries) GetLeaderboard(ctx context.Context) ([]GetLeaderboardRow, error) {
	rows, err := q.db.Query(ctx, getLeaderboard)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLeaderboardRow
	for rows.Next() {
		var i GetLeaderboardRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Points,
			&i.NrOfRounds,
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

const getPlayer = `-- name: GetPlayer :one
SELECT 
  p.id, 
  p.first_name, 
  p.last_name, 
  p.points, 
  p.hcp,
  u.email,
  u.id::UUID AS user_id
FROM players p LEFT JOIN users u ON u.id = p.user_id WHERE p.id = $1::UUID
`

type GetPlayerRow struct {
	ID        uuid.UUID       `json:"id"`
	FirstName string          `json:"firstName"`
	LastName  string          `json:"lastName"`
	Points    int32           `json:"points"`
	Hcp       decimal.Decimal `json:"hcp"`
	Email     string          `json:"email"`
	UserID    uuid.UUID       `json:"userId"`
}

func (q *Queries) GetPlayer(ctx context.Context, dollar_1 uuid.UUID) (GetPlayerRow, error) {
	row := q.db.QueryRow(ctx, getPlayer, dollar_1)
	var i GetPlayerRow
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Points,
		&i.Hcp,
		&i.Email,
		&i.UserID,
	)
	return i, err
}

const getPlayers = `-- name: GetPlayers :many
SELECT 
  p.id, 
  p.first_name, 
  p.last_name, 
  p.points, 
  p.hcp,
  u.email,
  u.id::UUID AS user_id
FROM players p LEFT JOIN users u ON u.id = p.user_id ORDER BY p.last_name ASC
`

type GetPlayersRow struct {
	ID        uuid.UUID       `json:"id"`
	FirstName string          `json:"firstName"`
	LastName  string          `json:"lastName"`
	Points    int32           `json:"points"`
	Hcp       decimal.Decimal `json:"hcp"`
	Email     string          `json:"email"`
	UserID    uuid.UUID       `json:"userId"`
}

func (q *Queries) GetPlayers(ctx context.Context) ([]GetPlayersRow, error) {
	rows, err := q.db.Query(ctx, getPlayers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPlayersRow
	for rows.Next() {
		var i GetPlayersRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Points,
			&i.Hcp,
			&i.Email,
			&i.UserID,
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

const updatePlayer = `-- name: UpdatePlayer :one
UPDATE players SET first_name = $1, last_name = $2, hcp = $3 WHERE id = $4 RETURNING id, first_name, last_name, points, user_id, hcp
`

type UpdatePlayerParams struct {
	FirstName string          `json:"firstName"`
	LastName  string          `json:"lastName"`
	Hcp       decimal.Decimal `json:"hcp"`
	ID        uuid.UUID       `json:"id"`
}

func (q *Queries) UpdatePlayer(ctx context.Context, arg *UpdatePlayerParams) (Player, error) {
	row := q.db.QueryRow(ctx, updatePlayer,
		arg.FirstName,
		arg.LastName,
		arg.Hcp,
		arg.ID,
	)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Points,
		&i.UserID,
		&i.Hcp,
	)
	return i, err
}
