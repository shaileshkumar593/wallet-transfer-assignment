package tests

import (
	"wallet-transfer-assignment/wallet-transfer/internal/repository"
	"wallet-transfer-assignment/wallet-transfer/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildTransferService(
	db *pgxpool.Pool,
) *service.TransferService {

	return service.NewTransferService(
		db,
		repository.NewWalletRepository(),
		repository.NewTransferRepository(),
		repository.NewLedgerRepository(),
		repository.NewIdempotencyRepository(),
	)
}
