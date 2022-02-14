package user

import (
	"context"
	"github.com/ddvkid/learngo/internal/storage"
	"github.com/ddvkid/learngo/internal/util"
)

type Request struct {
	EtherKey string `json:"ether_key" validate:"required"`
}

type Response struct {
	StarkKeys []string `json:"accounts"`
}

func List(ctx context.Context, store storage.Store, req Request) (*Response, error) {
	keys, err := store.ListStarkKeys(ctx, req.EtherKey)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, util.AccountNotFoundApiError(req.EtherKey)
	}

	return &Response{
		StarkKeys: keys,
	}, err
}
