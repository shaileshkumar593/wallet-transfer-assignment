package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type idempotencyRepository struct {
}

func NewIdempotencyRepository() IdempotencyRepository {
	return &idempotencyRepository{}
}

func (r *idempotencyRepository) Save(
	ctx context.Context,
	tx pgx.Tx,
	key string,
	transferID uuid.UUID,
	response []byte,
) error {

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO idempotency_records(
			idempotency_key,
			transfer_id,
			response
		)
		VALUES($1,$2,$3)
		`,
		key,
		transferID,
		response,
	)

	return err
}

func (r *idempotencyRepository) Get(
	ctx context.Context,
	tx pgx.Tx,
	key string,
) ([]byte, error) {

	var response []byte

	err := tx.QueryRow(
		ctx,
		`
		SELECT response
		FROM idempotency_records
		WHERE idempotency_key = $1
		`,
		key,
	).Scan(
		&response,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}
