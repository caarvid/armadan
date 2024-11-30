package database

import (
	"context"
	"fmt"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePool(ctx context.Context, host, port, name, user, password string) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable pool_max_conns=50",
		host,
		port,
		name,
		user,
		password,
	))

	if err != nil {
		return nil, err
	}

	dbConfig.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		pgxdecimal.Register(c.TypeMap())

		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)

	if err != nil {
		return nil, err
	}

	return pool, nil
}
