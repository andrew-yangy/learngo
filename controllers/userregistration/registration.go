package userregistration

import (
	"context"
	"github.com/ddvkid/learngo/internal/storage"
)

const signablePrefix = "Only sign this key linking request from Immutable X"

type Request struct {
	EtherKey       string `json:"ether_key" validate:"required"`
	StarkKey       string `json:"stark_key" validate:"required"`
	Nonce          uint64 `json:"nonce" validate:"bitlength31"`
	StarkSignature string `json:"stark_signature" validate:"required"`
}
type RequestVerifyEth struct {
	EtherKey       string `json:"ether_key" validate:"required"`
	StarkKey       string `json:"stark_key" validate:"required"`
	StarkSignature string `json:"stark_signature" validate:"required"`
	EthSignature   string `json:"eth_signature" validate:"required"`
}

type Response struct {
	TransactionHash string `json:"transaction_hash"`
}

type OperatorSignatureRequest struct {
	EtherKey string `json:"ether_key" validate:"required,eth_addr"`
	StarkKey string `json:"stark_key" validate:"required"`
}

type OperatorSignatureResponse struct {
	OperatorSignature string `json:"operator_signature"`
}

func Register(ctx context.Context, store storage.Store, req Request) (*Response, error) {
	account := &storage.Account{EtherKey: req.EtherKey, StarkKey: req.StarkKey, TxHash: ""}

	if err := Layer2Registration(ctx, store, account); err != nil {
		return nil, err
	}

	return &Response{TransactionHash: account.TxHash}, nil
}

type RegistrationBody struct {
	StarkKey string `json:"stark_key"`
	EtherKey string `json:"ether_key"`
}

func Layer2Registration(ctx context.Context, store storage.Store, account *storage.Account) error {
	return storage.InTx(ctx, store, func(tx storage.Tx) error {
		_, err := tx.InsertAccount(ctx, account)
		if err != nil {
			return err
		}

		// insert registration event for snapshot service to keep track of accounts
		//body, err := json.Marshal(RegistrationBody{StarkKey: account.StarkKey, EtherKey: account.EtherKey})
		//if err != nil {
		//	return err
		//}
		//_, err = tx.InsertEvent(ctx, &storage.Event{
		//	Type:    storage.RegistrationType,
		//	Status:  storage.AcceptedStatus, // placeholder - db complains if status not set
		//	RefID:   int64(id),
		//	BatchID: -1,
		//	Body:    body,
		//})
		//if err != nil {
		//	return err
		//}

		return nil
	})
}
