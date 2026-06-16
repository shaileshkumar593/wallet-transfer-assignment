package dto

type CreateTransferRequest struct {
	IdempotencyKey string `json:"idempotencyKey" binding:"required"`

	FromWalletID string `json:"fromWalletId" binding:"required"`

	ToWalletID string `json:"toWalletId" binding:"required"`

	Amount int64 `json:"amount" binding:"required,gt=0"`
}
