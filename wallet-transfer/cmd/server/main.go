package main

import (
	"log"

	"wallet-transfer-assignment/wallet-transfer/internal/config"
	"wallet-transfer-assignment/wallet-transfer/internal/database"
	"wallet-transfer-assignment/wallet-transfer/internal/handler"
	"wallet-transfer-assignment/wallet-transfer/internal/repository"
	"wallet-transfer-assignment/wallet-transfer/internal/router"
	"wallet-transfer-assignment/wallet-transfer/internal/service"
)

func main() {

	cfg := config.Load()

	dbPool, err :=
		database.NewPool(
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
		)

	if err != nil {
		log.Fatal(err)
	}

	defer dbPool.Close()

	//--------------------------------------------------
	// Repositories
	//--------------------------------------------------

	walletRepo :=
		repository.NewWalletRepository()

	transferRepo :=
		repository.NewTransferRepository()

	ledgerRepo :=
		repository.NewLedgerRepository()

	idempotencyRepo :=
		repository.NewIdempotencyRepository()

	//--------------------------------------------------
	// Services
	//--------------------------------------------------

	transferService :=
		service.NewTransferService(
			dbPool,
			walletRepo,
			transferRepo,
			ledgerRepo,
			idempotencyRepo,
		)

	walletService :=
		service.NewWalletService(
			dbPool,
			walletRepo,
		)

	//--------------------------------------------------
	// Handlers
	//--------------------------------------------------

	transferHandler :=
		handler.NewTransferHandler(
			transferService,
		)

	walletHandler :=
		handler.NewWalletHandler(
			walletService,
		)

	healthHandler :=
		handler.NewHealthHandler()

	//--------------------------------------------------
	// Router
	//--------------------------------------------------

	r := router.Setup(
		transferHandler,
		walletHandler,
		healthHandler,
	)

	log.Printf(
		"server running on :%s",
		cfg.ServerPort,
	)

	if err := r.Run(
		":" + cfg.ServerPort,
	); err != nil {

		log.Fatal(err)
	}
}
