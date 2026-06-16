package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransferStatus string

const (
	Pending   TransferStatus = "PENDING"
	Processed TransferStatus = "PROCESSED"
	Failed    TransferStatus = "FAILED"
)

type Transfer struct {
	ID             uuid.UUID
	IdempotencyKey string

	FromWalletID uuid.UUID
	ToWalletID   uuid.UUID

	Amount int64

	Status TransferStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Transfer) Process() error {

	if t.Status != Pending {
		return ErrInvalidState
	}

	t.Status = Processed

	return nil
}

func (t *Transfer) Fail() error {

	if t.Status != Pending {
		return ErrInvalidState
	}

	t.Status = Failed

	return nil
}
