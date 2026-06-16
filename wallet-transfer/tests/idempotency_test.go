package tests

import (
	"context"
	"sync"
	"testing"
	"wallet-transfer-assignment/wallet-transfer/internal/dto"

	"github.com/stretchr/testify/assert"
)

func TestIdempotency(
	t *testing.T,
) {

	ctx := context.Background()

	db := SetupDB(t)

	_ = Seed(
		ctx,
		db,
	)

	service :=
		BuildTransferService(
			db,
		)

	req :=
		dto.CreateTransferRequest{
			IdempotencyKey: "same-key",

			FromWalletID: "11111111-1111-1111-1111-111111111111",

			ToWalletID: "22222222-2222-2222-2222-222222222222",

			Amount: 100,
		}

	wg := sync.WaitGroup{}

	for i := 0; i < 50; i++ {

		wg.Add(1)

		go func() {

			defer wg.Done()

			_, _ =
				service.Transfer(
					ctx,
					req,
				)

		}()
	}

	wg.Wait()

	var transfers int

	err := db.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM transfers
		`,
	).Scan(
		&transfers,
	)

	assert.NoError(
		t,
		err,
	)

	assert.Equal(
		t,
		1,
		transfers,
	)

	var ledgers int

	err = db.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM ledger_entries
		`,
	).Scan(
		&ledgers,
	)

	assert.NoError(
		t,
		err,
	)

	assert.Equal(
		t,
		2,
		ledgers,
	)
}
