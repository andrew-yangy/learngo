package storage

import (
	"encoding/json"
	"time"
)

type TransactionType string

const (
	DepositType             = TransactionType("deposit")
	MintType                = TransactionType("mint")
	WithdrawalType          = TransactionType("withdrawal")
	FullWithdrawalType      = TransactionType("full_withdrawal")
	FalseFullWithdrawalType = TransactionType("false_full_withdrawal")
	TransferType            = TransactionType("transfer")
	SettlementType          = TransactionType("settlement")
	OrderType               = TransactionType("order")
	RegistrationType        = TransactionType("registration")
	RootUpdateType          = TransactionType("root_update")
)

type Status string

const (
	CreatedStatus         = Status("created")
	PendingStatus         = Status("pending")
	AcceptedStatus        = Status("accepted")
	RejectedStatus        = Status("rejected")
	ConfirmedStatus       = Status("confirmed")
	RolledBackStatus      = Status("rolled_back")
	CancelledStatus       = Status("cancelled")
	FilledStatus          = Status("filled")
	PartiallyFilledStatus = Status("partially_filled")
	ExpiredStatus         = Status("expired")
	WithdrawnStatus       = Status("withdrawn")
)

type Account struct {
	ID       uint64 `db:"id"`
	StarkKey string `db:"stark_key"`
	EtherKey string `db:"ether_key"`
	Nonce    uint64 `db:"nonce"`
	TxHash   string `db:"tx_hash"`
}

type Event struct {
	ID        uint64          `db:"id"`
	Type      TransactionType `db:"type"`
	Status    Status          `db:"status"`
	RefID     int64           `db:"reference_id"`
	BatchID   int64           `db:"batch_id"`
	Timestamp time.Time       `db:"creation_time"`
	Reason    string          `db:"reason"`
	Context   string          `db:"context"`
	Body      json.RawMessage `db:"body"`
}
