package service

import (
	"context"

	"wallet-transfer-assignment/wallet-transfer/internal/domain"
	"wallet-transfer-assignment/wallet-transfer/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletService struct {
	db *pgxpool.Pool

	walletRepo repository.WalletRepository
}

func NewWalletService(
	db *pgxpool.Pool,
	walletRepo repository.WalletRepository,
) *WalletService {

	return &WalletService{
		db:         db,
		walletRepo: walletRepo,
	}
}

func (s *WalletService) GetWallet(
	ctx context.Context,
	id string,
) (*domain.Wallet, error) {

	tx, err := s.db.Begin(ctx)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	walletID, err :=
		uuid.Parse(id)

	if err != nil {
		return nil, err
	}

	wallet, err :=
		s.walletRepo.GetByID(
			ctx,
			tx,
			walletID,
		)

	if err != nil {
		return nil, err
	}

	if err :=
		tx.Commit(ctx); err != nil {

		return nil, err
	}

	return wallet, nil
}
