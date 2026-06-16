package router

import (
	"wallet-transfer-assignment/wallet-transfer/internal/handler"
	"wallet-transfer-assignment/wallet-transfer/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(
	transferHandler *handler.TransferHandler,
	walletHandler *handler.WalletHandler,
	healthHandler *handler.HealthHandler,
) *gin.Engine {

	r := gin.New()

	r.Use(
		middleware.Logger(),
	)

	r.Use(
		middleware.RequestID(),
	)

	r.Use(
		gin.Recovery(),
	)

	r.GET(
		"/health",
		healthHandler.Health,
	)

	r.POST(
		"/transfers",
		transferHandler.CreateTransfer,
	)

	r.GET(
		"/wallets/:id",
		walletHandler.GetWallet,
	)

	return r
}
