package service

import (
	"context"
	"encoding/json"
	"errors"

	"wallet-transfer-assignment/wallet-transfer/internal/domain"
	"wallet-transfer-assignment/wallet-transfer/internal/dto"
	"wallet-transfer-assignment/wallet-transfer/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransferService struct {
	db *pgxpool.Pool

	walletRepo      repository.WalletRepository
	transferRepo    repository.TransferRepository
	ledgerRepo      repository.LedgerRepository
	idempotencyRepo repository.IdempotencyRepository
}

func NewTransferService(
	db *pgxpool.Pool,

	walletRepo repository.WalletRepository,
	transferRepo repository.TransferRepository,
	ledgerRepo repository.LedgerRepository,
	idempotencyRepo repository.IdempotencyRepository,
) *TransferService {

	return &TransferService{
		db: db,

		walletRepo:      walletRepo,
		transferRepo:    transferRepo,
		ledgerRepo:      ledgerRepo,
		idempotencyRepo: idempotencyRepo,
	}
}

func (s *TransferService) Transfer(
	ctx context.Context,
	req dto.CreateTransferRequest,
) (*dto.TransferResponse, error) {

	tx, err := s.db.BeginTx(
		ctx,
		pgx.TxOptions{
			IsoLevel: pgx.Serializable,
		},
	)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	//---------------------------------------------------------
	// STEP 1
	// IDEMPOTENCY CHECK
	//---------------------------------------------------------

	existingResponse, err :=
		s.idempotencyRepo.Get(
			ctx,
			tx,
			req.IdempotencyKey,
		)

	if err == nil {

		var response dto.TransferResponse

		if err := json.Unmarshal(
			existingResponse,
			&response,
		); err != nil {
			return nil, err
		}

		return &response, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	//---------------------------------------------------------
	// STEP 2
	// VALIDATE UUIDS
	//---------------------------------------------------------

	fromWalletID, err :=
		uuid.Parse(req.FromWalletID)

	if err != nil {
		return nil, err
	}

	toWalletID, err :=
		uuid.Parse(req.ToWalletID)

	if err != nil {
		return nil, err
	}

	//---------------------------------------------------------
	// STEP 3
	// LOCK WALLETS
	//---------------------------------------------------------

	wallets, err :=
		s.walletRepo.LockWallets(
			ctx,
			tx,
			fromWalletID,
			toWalletID,
		)

	if err != nil {
		return nil, err
	}

	fromWallet, ok :=
		wallets[fromWalletID]

	if !ok {
		return nil, domain.ErrWalletNotFound
	}

	toWallet, ok :=
		wallets[toWalletID]

	if !ok {
		return nil, domain.ErrWalletNotFound
	}

	//---------------------------------------------------------
	// STEP 4
	// CREATE TRANSFER (PENDING)
	//---------------------------------------------------------

	transfer := &domain.Transfer{
		ID:             uuid.New(),
		IdempotencyKey: req.IdempotencyKey,
		FromWalletID:   fromWalletID,
		ToWalletID:     toWalletID,
		Amount:         req.Amount,
		Status:         domain.Pending,
	}

	err = s.transferRepo.Create(
		ctx,
		tx,
		transfer,
	)

	if err != nil {
		return nil, err
	}

	//---------------------------------------------------------
	// STEP 5
	// BALANCE CHECK
	//---------------------------------------------------------

	if fromWallet.Balance < req.Amount {

		_ = s.transferRepo.MarkFailed(
			ctx,
			tx,
			transfer.ID,
		)

		return nil,
			domain.ErrInsufficientFunds
	}

	//---------------------------------------------------------
	// STEP 6
	// UPDATE BALANCES
	//---------------------------------------------------------

	newFromBalance :=
		fromWallet.Balance -
			req.Amount

	newToBalance :=
		toWallet.Balance +
			req.Amount

	err = s.walletRepo.UpdateBalance(
		ctx,
		tx,
		fromWallet.ID,
		newFromBalance,
	)

	if err != nil {

		_ = s.transferRepo.MarkFailed(
			ctx,
			tx,
			transfer.ID,
		)

		return nil, err
	}

	err = s.walletRepo.UpdateBalance(
		ctx,
		tx,
		toWallet.ID,
		newToBalance,
	)

	if err != nil {

		_ = s.transferRepo.MarkFailed(
			ctx,
			tx,
			transfer.ID,
		)

		return nil, err
	}

	//---------------------------------------------------------
	// STEP 7
	// CREATE LEDGER ENTRIES
	//---------------------------------------------------------

	debit := &domain.LedgerEntry{
		ID:         uuid.New(),
		TransferID: transfer.ID,
		WalletID:   fromWallet.ID,
		Type:       domain.Debit,
		Amount:     req.Amount,
	}

	credit := &domain.LedgerEntry{
		ID:         uuid.New(),
		TransferID: transfer.ID,
		WalletID:   toWallet.ID,
		Type:       domain.Credit,
		Amount:     req.Amount,
	}

	err = s.ledgerRepo.CreateEntries(
		ctx,
		tx,
		debit,
		credit,
	)

	if err != nil {

		_ = s.transferRepo.MarkFailed(
			ctx,
			tx,
			transfer.ID,
		)

		return nil, err
	}

	//---------------------------------------------------------
	// STEP 8
	// MARK PROCESSED
	//---------------------------------------------------------

	err = s.transferRepo.UpdateStatus(
		ctx,
		tx,
		transfer.ID,
		domain.Processed,
	)

	if err != nil {
		return nil, err
	}

	//---------------------------------------------------------
	// STEP 9
	// BUILD RESPONSE
	//---------------------------------------------------------

	response := dto.TransferResponse{
		TransferID: transfer.ID.String(),
		Status:     string(domain.Processed),
		Amount:     req.Amount,
	}

	responseBytes, err :=
		json.Marshal(response)

	if err != nil {
		return nil, err
	}

	//---------------------------------------------------------
	// STEP 10
	// SAVE IDEMPOTENCY RECORD
	//---------------------------------------------------------

	err = s.idempotencyRepo.Save(
		ctx,
		tx,
		req.IdempotencyKey,
		transfer.ID,
		responseBytes,
	)

	if err != nil {
		return nil, err
	}

	//---------------------------------------------------------
	// STEP 11
	// COMMIT
	//---------------------------------------------------------

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &response, nil
}
