package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ddvkid/learngo/internal/storage"
)

type AccountAlreadyExists struct{}

func (e *AccountAlreadyExists) Error() string {
	return "account already exists"
}

func (s Source) InsertAccount(ctx context.Context, account *storage.Account) (uint64, error) {
	var id uint64
	query := `
	INSERT INTO accounts (stark_key, ether_key, nonce, tx_hash)
	VALUES ($1, $2, $3, $4) ON CONFLICT (stark_key, ether_key) DO NOTHING
	RETURNING id
	`
	if err := s.GetContext(ctx, &id, query, account.StarkKey, account.EtherKey, account.Nonce, account.TxHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, &AccountAlreadyExists{}
		}
		return 0, err
	}
	return id, nil
}

func (s Source) ListStarkKeys(ctx context.Context, etherKey string) ([]string, error) {
	var starkKeys []string
	query := `SELECT stark_key FROM accounts WHERE ether_key = $1`
	if err := s.SelectContext(ctx, &starkKeys, query, etherKey); err != nil {
		return nil, err
	}
	return starkKeys, nil
}
