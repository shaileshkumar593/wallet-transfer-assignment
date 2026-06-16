package tests

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func SetupDB(
	t *testing.T,
) *pgxpool.Pool {

	ctx := context.Background()

	container, err :=
		postgres.Run(
			ctx,
			"postgres:17",
			postgres.WithDatabase("wallet"),
			postgres.WithUsername("wallet"),
			postgres.WithPassword("wallet"),
		)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = container.Terminate(ctx)
	})

	connString, err :=
		container.ConnectionString(
			ctx,
			"sslmode=disable",
		)

	if err != nil {
		t.Fatal(err)
	}

	db, err :=
		pgxpool.New(
			ctx,
			connString,
		)

	if err != nil {
		t.Fatal(err)
	}

	migration, err :=
		os.ReadFile(
			filepath.Join(
				"..",
				"migrations",
				"001_init.up.sql",
			),
		)

	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(
		ctx,
		string(migration),
	)

	if err != nil {
		t.Fatal(err)
	}

	return db
}
