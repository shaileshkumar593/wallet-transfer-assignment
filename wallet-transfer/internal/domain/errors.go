package domain

import "errors"

var (
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrTransferNotFound  = errors.New("transfer not found")
	ErrDuplicateTransfer = errors.New("duplicate transfer")
	ErrInvalidState      = errors.New("invalid state transition")
)
