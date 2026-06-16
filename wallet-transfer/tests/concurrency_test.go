package tests

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"wallet-transfer-assignment/wallet-transfer/internal/dto"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentDebit(
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

	var success int32

	wg := sync.WaitGroup{}

	for i := 0; i < 2; i++ {

		wg.Add(1)

		go func(
			index int,
		) {

			defer wg.Done()

			req :=
				dto.CreateTransferRequest{
					IdempotencyKey: fmt.Sprintf(
						"key-%d",
						index,
					),

					FromWalletID: "11111111-1111-1111-1111-111111111111",

					ToWalletID: "22222222-2222-2222-2222-222222222222",

					Amount: 100,
				}

			_, err :=
				service.Transfer(
					ctx,
					req,
				)

			if err == nil {

				atomic.AddInt32(
					&success,
					1,
				)
			}

		}(i)
	}

	wg.Wait()

	assert.Equal(
		t,
		int32(1),
		success,
	)

	var balance int64

	err := db.QueryRow(
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
