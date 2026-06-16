package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(
	host string,
	port string,
	user string,
	password string,
	db string,
) (*pgxpool.Pool, error) {

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		db,
	)

	cfg, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 20
	cfg.MinConns = 5

	pool, err :=
		pgxpool.NewWithConfig(
			context.Background(),
			cfg,
		)

	if err != nil {
		return nil, err
	}

	if err := pool.Ping(
		context.Background(),
	); err != nil {

		return nil, err
	}

	return pool, nil
}
