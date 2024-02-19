package database

import (
	"context"

	"github.com/glu-project/pkg/postgres"

	"github.com/jackc/pgx/v5"
)

type Txer interface {
	Transaction(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error
}

type txer struct {
	client *postgres.PostgresClient
}

func NewTxer(client *postgres.PostgresClient) Txer {
	return &txer{
		client: client,
	}
}

func (t *txer) Transaction(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := t.client.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := fn(ctx, tx); err != nil {
		return err
	}

	tx.Commit(ctx)

	return nil
}
