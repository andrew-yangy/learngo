package storage

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Source interface {
	InsertAccount(ctx context.Context, account *Account) (uint64, error)
	ListStarkKeys(ctx context.Context, etherKey string) ([]string, error)
	InsertEvent(ctx context.Context, e *Event) (uint64, error)
}

type Store interface {
	Source
	BeginTx(ctx context.Context) (Tx, error)
}

type Tx interface {
	Source
	Commit() error
	Rollback() error
}

func InTx(ctx context.Context, store Store, fn func(tx Tx) error) error {
	tx, err := store.BeginTx(ctx)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		logrus.Printf("error in transaction: %+v", err)
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	} else {
		return tx.Commit()
	}
}
