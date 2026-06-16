package repository

import (
	"context"

	"wallet-transfer-assignment/wallet-transfer/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type transferRepository struct {
}

func NewTransferRepository() TransferRepository {
	return &transferRepository{}
}

func (r *transferRepository) Create(
	ctx context.Context,
	tx pgx.Tx,
	transfer *domain.Transfer,
) error {

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO transfers(
			id,
			idempotency_key,
			from_wallet_id,
			to_wallet_id,
			amount,
			status
		)
		VALUES($1,$2,$3,$4,$5,$6)
		`,
		transfer.ID,
		transfer.IdempotencyKey,
		transfer.FromWalletID,
		transfer.ToWalletID,
		transfer.Amount,
		transfer.Status,
	)

	return err
}

func (r *transferRepository) UpdateStatus(
	ctx context.Context,
	tx pgx.Tx,
	transferID uuid.UUID,
	status domain.TransferStatus,
) error {

	_, err := tx.Exec(
		ctx,
		`
		UPDATE transfers
		SET
			status = $1,
			updated_at = NOW()
		WHERE id = $2
		`,
		status,
		transferID,
	)

	return err
}

func (r *transferRepository) GetByIdempotencyKey(
	ctx context.Context,
	tx pgx.Tx,
	key string,
) (*domain.Transfer, error) {

	var transfer domain.Transfer

	err := tx.QueryRow(
		ctx,
		`
		SELECT
			id,
			idempotency_key,
			from_wallet_id,
			to_wallet_id,
			amount,
			status,
			created_at,
			updated_at
		FROM transfers
		WHERE idempotency_key = $1
		`,
		key,
	).Scan(
		&transfer.ID,
		&transfer.IdempotencyKey,
		&transfer.FromWalletID,
		&transfer.ToWalletID,
		&transfer.Amount,
		&transfer.Status,
		&transfer.CreatedAt,
		&transfer.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &transfer, nil
}

func (r *transferRepository) MarkFailed(
	ctx context.Context,
	tx pgx.Tx,
	transferID uuid.UUID,
) error {

	_, err := tx.Exec(
		ctx,
		`
		UPDATE transfers
		SET status = 'FAILED',
		    updated_at = NOW()
		WHERE id = $1
		`,
		transferID,
	)

	return err
}
