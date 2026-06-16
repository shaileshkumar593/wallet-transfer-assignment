package tests

import (
	"context"
	"testing"
	"wallet-transfer-assignment/wallet-transfer/internal/dto"

	"github.com/stretchr/testify/assert"
)

func TestTransferSuccess(
	t *testing.T,
) {

	ctx := context.Background()

	db := SetupDB(t)

	err := Seed(
		ctx,
		db,
	)

	assert.NoError(
		t,
		err,
	)

	service := BuildTransferService(
		db,
	)

	req :=
		dto.CreateTransferRequest{
			IdempotencyKey: "transfer-1",

			FromWalletID: "11111111-1111-1111-1111-111111111111",

			ToWalletID: "22222222-2222-2222-2222-222222222222",

			Amount: 100,
		}

	resp, err :=
		service.Transfer(
			ctx,
			req,
		)

	assert.NoError(
		t,
		err,
	)

	assert.Equal(
		t,
		"PROCESSED",
		resp.Status,
	)

	var balance int64

	err = db.QueryRow(
		ctx,
		`
		SELECT balance
		FROM wallets
		WHERE id='11111111-1111-1111-1111-111111111111'
		`,
	).Scan(
		&balance,
	)

	assert.NoError(
		t,
		err,
	)

	assert.Equal(
		t,
		int64(0),
		balance,
	)
}
