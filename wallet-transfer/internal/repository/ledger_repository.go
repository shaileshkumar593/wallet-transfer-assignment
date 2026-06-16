package repository

import (
	"context"

	"wallet-transfer-assignment/wallet-transfer/internal/domain"

	"github.com/jackc/pgx/v5"
)

type ledgerRepository struct {
}

func NewLedgerRepository() LedgerRepository {
	return &ledgerRepository{}
}

func (r *ledgerRepository) CreateEntries(
	ctx context.Context,
	tx pgx.Tx,
	debit *domain.LedgerEntry,
	credit *domain.LedgerEntry,
) error {

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO ledger_entries(
			id,
			transfer_id,
			wallet_id,
			entry_type,
			amount
		)
		VALUES
		($1,$2,$3,$4,$5),
		($6,$7,$8,$9,$10)
		`,
		debit.ID,
		debit.TransferID,
		debit.WalletID,
		debit.Type,
		debit.Amount,

		credit.ID,
		credit.TransferID,
		credit.WalletID,
		credit.Type,
		credit.Amount,
	)

	return err
}
