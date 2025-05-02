package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/caarvid/armadan/internal/utils"
	"go.step.sm/crypto/randutil"
)

func toPlayer(entity any) *armadan.Player {
	switch p := entity.(type) {
	case schema.PlayersExtended:
		return &armadan.Player{
			ID:        p.ID,
			FirstName: p.FirstName,
			LastName:  p.LastName,
			Points:    p.Points,
			UserID:    p.UserID,
			Email:     p.Email,
			Hcp:       p.Hcp,
		}
	case schema.Player:
		return &armadan.Player{
			ID:        p.ID,
			FirstName: p.FirstName,
			LastName:  p.LastName,
			UserID:    p.UserID,
		}
	case schema.GetRemainingPlayersByResultIdRow:
		return &armadan.Player{
			ID:        p.ID,
			FirstName: p.FirstName,
			LastName:  p.LastName,
			UserID:    p.UserID,
			Hcp:       p.Hcp,
		}
	}

	return &armadan.Player{}
}

type players struct {
	dbReader schema.Querier
	dbWriter schema.Querier
	pool     *sql.DB
}

func NewPlayerService(reader, writer schema.Querier, pool *sql.DB) *players {
	return &players{
		dbReader: reader,
		dbWriter: writer,
		pool:     pool,
	}
}

func (ps *players) All(ctx context.Context) ([]armadan.Player, error) {
	p, err := ps.dbReader.GetPlayers(ctx)

	if err != nil {
		return nil, err
	}

	return armadan.MapEntities(p, toPlayer), nil
}

func (ps *players) Get(ctx context.Context, id string) (*armadan.Player, error) {
	p, err := ps.dbReader.GetPlayer(ctx, id)

	if err != nil {
		return nil, err
	}

	return toPlayer(p), nil
}

func (ps *players) GetByUserId(ctx context.Context, id string) (*armadan.Player, error) {
	p, err := ps.dbReader.GetPlayerByUserId(ctx, id)
	if err != nil {
		return nil, err
	}

	return toPlayer(p), nil
}

func (ps *players) Create(ctx context.Context, data *armadan.Player) (*armadan.Player, error) {
	pw, err := randutil.Alphanumeric(24)
	if err != nil {
		return nil, err
	}

	hash, err := utils.GenerateHash(pw, nil)
	if err != nil {
		return nil, err
	}

	tx, err := ps.pool.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()
	qtx := schema.New(tx)

	user, err := qtx.CreateUser(ctx, &schema.CreateUserParams{
		ID:       armadan.GetId(),
		Email:    data.Email,
		Password: hash.Encode(),
	})

	if err != nil {
		return nil, err
	}

	player, err := qtx.CreatePlayer(ctx, &schema.CreatePlayerParams{
		ID:        armadan.GetId(),
		FirstName: data.FirstName,
		LastName:  data.LastName,
		UserID:    user.ID,
	})

	if err != nil {
		return nil, err
	}

	_, err = qtx.CreateHcpChange(ctx, &schema.CreateHcpChangeParams{
		OldHcp:    data.Hcp,
		NewHcp:    data.Hcp,
		ValidFrom: time.Now().Format(armadan.DEFAULT_TIME_FORMAT),
		PlayerID:  sql.NullString{String: player.ID, Valid: true},
	})

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &armadan.Player{
		ID:        player.ID,
		FirstName: player.FirstName,
		LastName:  player.LastName,
		UserID:    player.UserID,
		Email:     user.Email,
		Hcp:       data.Hcp,
	}, nil
}

func (ps *players) Update(ctx context.Context, data *armadan.Player) (*armadan.Player, error) {
	tx, err := ps.pool.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()
	qtx := schema.New(tx)

	player, err := qtx.UpdatePlayer(ctx, &schema.UpdatePlayerParams{
		ID:        data.ID,
		FirstName: data.FirstName,
		LastName:  data.LastName,
	})

	if err != nil {
		return nil, err
	}

	user, err := qtx.UpdateUserEmail(ctx, &schema.UpdateUserEmailParams{
		Email: data.Email,
		ID:    player.UserID,
	})

	if err != nil {
		return nil, err
	}

	currentHcp, _ := qtx.GetPlayerHcp(ctx, sql.NullString{String: player.ID, Valid: true})

	if currentHcp != data.Hcp {
		_, err = qtx.CreateHcpChange(ctx, &schema.CreateHcpChangeParams{
			OldHcp:    currentHcp,
			NewHcp:    data.Hcp,
			ValidFrom: time.Now().Format(armadan.DEFAULT_TIME_FORMAT),
			PlayerID:  sql.NullString{String: player.ID, Valid: true},
		})

		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &armadan.Player{
		ID:        player.ID,
		FirstName: player.FirstName,
		LastName:  player.LastName,
		UserID:    player.UserID,
		Email:     user.Email,
		Hcp:       data.Hcp,
	}, nil
}

func (ps *players) Delete(ctx context.Context, id string) error {
	return ps.dbWriter.DeleteUser(ctx, id)
}
