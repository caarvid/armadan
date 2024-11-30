package service

import (
	"context"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/caarvid/armadan/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.step.sm/crypto/randutil"
)

func toPlayer(entity any) *armadan.Player {
	switch p := entity.(type) {
	case schema.GetPlayerRow:
		return &armadan.Player{
			ID:        p.ID,
			FirstName: p.FirstName,
			LastName:  p.LastName,
			Points:    p.Points,
			UserID:    p.UserID,
			Email:     p.Email,
			Hcp:       p.Hcp,
		}
	case schema.GetPlayersRow:
		return &armadan.Player{
			ID:        p.ID,
			FirstName: p.FirstName,
			LastName:  p.LastName,
			Points:    p.Points,
			UserID:    p.UserID,
			Email:     p.Email,
			Hcp:       p.Hcp,
		}
	}

	return &armadan.Player{}
}

type players struct {
	db   schema.Querier
	pool *pgxpool.Pool
}

func NewPlayerService(db schema.Querier, pool *pgxpool.Pool) *players {
	return &players{
		db:   db,
		pool: pool,
	}
}

func (ps *players) All(ctx context.Context) ([]armadan.Player, error) {
	p, err := ps.db.GetPlayers(ctx)

	if err != nil {
		return nil, err
	}

	return armadan.MapEntities(p, toPlayer), nil
}

func (ps *players) Get(ctx context.Context, id uuid.UUID) (*armadan.Player, error) {
	p, err := ps.db.GetPlayer(ctx, id)

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

	tx, err := ps.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)
	qtx := schema.New(tx)

	user, err := qtx.CreateUser(ctx, &schema.CreateUserParams{
		Email:    data.Email,
		Password: hash.Encode(),
	})

	if err != nil {
		return nil, err
	}

	player, err := qtx.CreatePlayer(ctx, &schema.CreatePlayerParams{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		UserID:    user.ID,
		Hcp: pgtype.Numeric{
			Int:              data.Hcp.BigInt(),
			Exp:              data.Hcp.Exponent(),
			InfinityModifier: pgtype.Finite,
			NaN:              false,
			Valid:            true,
		},
	})

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &armadan.Player{
		ID:        player.ID,
		FirstName: player.FirstName,
		LastName:  player.LastName,
		UserID:    player.UserID,
		Email:     user.Email,
		Hcp:       player.Hcp,
	}, nil
}

func (ps *players) Update(ctx context.Context, data *armadan.Player) (*armadan.Player, error) {
	tx, err := ps.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)
	qtx := schema.New(tx)

	player, err := qtx.UpdatePlayer(ctx, &schema.UpdatePlayerParams{
		ID:        data.ID,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Hcp: pgtype.Numeric{
			Int:              data.Hcp.BigInt(),
			Exp:              data.Hcp.Exponent(),
			InfinityModifier: pgtype.Finite,
			NaN:              false,
			Valid:            true,
		},
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

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &armadan.Player{
		ID:        player.ID,
		FirstName: player.FirstName,
		LastName:  player.LastName,
		UserID:    player.UserID,
		Email:     user.Email,
		Hcp:       player.Hcp,
	}, nil
}

func (ps *players) Delete(ctx context.Context, id uuid.UUID) error {
	return ps.db.DeletePlayer(ctx, id)
}
