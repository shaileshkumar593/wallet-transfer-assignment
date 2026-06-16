package repository

import (
	"context"

	"wallet-transfer-assignment/wallet-transfer/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type walletRepository struct {
}

func NewWalletRepository() WalletRepository {
	return &walletRepository{}
}

func (r *walletRepository) LockWallets(
	ctx context.Context,
	tx pgx.Tx,
	fromID uuid.UUID,
	toID uuid.UUID,
) (map[uuid.UUID]domain.Wallet, error) {

	rows, err := tx.Query(
		ctx,
		`
		SELECT
			id,
			balance,
			created_at
		FROM wallets
		WHERE id IN ($1,$2)
		ORDER BY id
		FOR UPDATE
		`,
		fromID,
		toID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	wallets := make(
		map[uuid.UUID]domain.Wallet,
	)

	for rows.Next() {

		var wallet domain.Wallet

		err := rows.Scan(
			&wallet.ID,
			&wallet.Balance,
			&wallet.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		wallets[wallet.ID] = wallet
	}

	return wallets, nil
}

func (r *walletRepository) UpdateBalance(
	ctx context.Context,
	tx pgx.Tx,
	walletID uuid.UUID,
	balance int64,
) error {

	_, err := tx.Exec(
		ctx,
		`
		UPDATE wallets
		SET balance = $1
		WHERE id = $2
		`,
		balance,
		walletID,
	)

	return err
}

func (r *walletRepository) GetByID(
	ctx context.Context,
	tx pgx.Tx,
	id uuid.UUID,
) (*domain.Wallet, error) {

	var wallet domain.Wallet

	err := tx.QueryRow(
		ctx,
		`
		SELECT
			id,
			balance,
			created_at
		FROM wallets
		WHERE id = $1
		`,
		id,
	).Scan(
		&wallet.ID,
		&wallet.Balance,
		&wallet.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}
