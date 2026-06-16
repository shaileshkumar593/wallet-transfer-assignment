package repository

import (
	"context"

	"wallet-transfer-assignment/wallet-transfer/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type WalletRepository interface {
	LockWallets(
		ctx context.Context,
		tx pgx.Tx,
		fromID uuid.UUID,
		toID uuid.UUID,
	) (map[uuid.UUID]domain.Wallet, error)

	UpdateBalance(
		ctx context.Context,
		tx pgx.Tx,
		walletID uuid.UUID,
		balance int64,
	) error

	GetByID(
		ctx context.Context,
		tx pgx.Tx,
		id uuid.UUID,
	) (*domain.Wallet, error)
}

type TransferRepository interface {
	Create(
		ctx context.Context,
		tx pgx.Tx,
		transfer *domain.Transfer,
	) error

	UpdateStatus(
		ctx context.Context,
		tx pgx.Tx,
		transferID uuid.UUID,
		status domain.TransferStatus,
	) error

	GetByIdempotencyKey(
		ctx context.Context,
		tx pgx.Tx,
		key string,
	) (*domain.Transfer, error)

	MarkFailed(
		ctx context.Context,
		tx pgx.Tx,
		transferID uuid.UUID,
	) error
}

type LedgerRepository interface {
	CreateEntries(
		ctx context.Context,
		tx pgx.Tx,
		debit *domain.LedgerEntry,
		credit *domain.LedgerEntry,
	) error
}

type IdempotencyRepository interface {
	Save(
		ctx context.Context,
		tx pgx.Tx,
		key string,
		transferID uuid.UUID,
		response []byte,
	) error

	Get(
		ctx context.Context,
		tx pgx.Tx,
		key string,
	) ([]byte, error)
}
