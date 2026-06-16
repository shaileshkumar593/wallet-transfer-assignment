package domain

import (
	"time"

	"github.com/google/uuid"
)

type LedgerType string

const (
	Debit  LedgerType = "DEBIT"
	Credit LedgerType = "CREDIT"
)

type LedgerEntry struct {
	ID         uuid.UUID
	TransferID uuid.UUID
	WalletID   uuid.UUID

	Type   LedgerType
	Amount int64

	CreatedAt time.Time
}
