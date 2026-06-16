package tests

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Seed(
	ctx context.Context,
	db *pgxpool.Pool,
) error {

	sql := `
TRUNCATE idempotency_records CASCADE;
TRUNCATE ledger_entries CASCADE;
TRUNCATE transfers CASCADE;
TRUNCATE wallets CASCADE;

INSERT INTO wallets(id,balance)
VALUES
('11111111-1111-1111-1111-111111111111',100),
('22222222-2222-2222-2222-222222222222',0);
`

	_, err := db.Exec(
		ctx,
		sql,
	)

	return err
}
