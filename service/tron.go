package service

import (
	"context"

	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
)

type Tron interface {
	Close()
	GetWallet(ctx context.Context, seedPhrase *string) *model.TrxHdWallet
	// GetWallet(ctx context.Context, passphrase, file string) (*keystore.Key, error)
	GetBalance(ctx context.Context, address *string) (balance *int64, err error)
	SendTx(ctx context.Context, txOpts *model.TxOpts, wallet *model.TrxHdWallet) (transaction *api.TransactionExtention, err error)
	GetTx(ctx context.Context, txhash string) (*core.TransactionInfo, error)
	GetUnconfirmedTxAddress(ctx context.Context, address *string) (*GetTransactionResponse, error)
	GetConfirmedTxAddress(ctx context.Context, address *string) (*GetTransactionResponse, error)
	CheckAddress(address string) error
	GetCurrentBlock(ctx context.Context) (*int64, error)
	GetTxByAccountAddress(ctx context.Context, address *string, filter *GetTxByAccountAddressFilter) (*GetTransactionResponse, error)
}

type GetTransactionResponse struct {
	Data    []TransactionData `json:"data"`
	Success bool              `json:"success"`
	Meta    *Meta             `json:"meta"`
}

type TransactionData struct {
	Ret []struct {
		ContractRet *string `json:"contractRet"`
		Fee         *int64  `json:"fee"`
	} `json:"ret"`
	Signature            []string            `json:"signature"`
	TxID                 *string             `json:"txID"`
	NetUsage             *int64              `json:"net_usage"`
	RawDataHex           *string             `json:"raw_data_hex"`
	NetFee               *int64              `json:"net_fee"`
	EnergyUsage          *int64              `json:"energy_usage"`
	BlockNumber          *int64              `json:"blockNumber"`
	BlockTimestamp       *int64              `json:"block_timestamp"`
	EnergyFee            *int64              `json:"energy_fee"`
	EnergyUsageTotal     *int64              `json:"energy_usage_total"`
	RawData              *RawDataTransaction `json:"raw_data"`
	InternalTransactions []interface{}       `json:"internal_transactions"`
}

type RawDataTransaction struct {
	Contract      []Contract `json:"contract"`
	RefBlockBytes *string    `json:"ref_block_bytes"`
	RefBlockHash  *string    `json:"ref_block_hash"`
	Expiration    *int64     `json:"expiration"`
	Timestamp     *int64     `json:"timestamp"`
}

type Contract struct {
	Parameter struct {
		Value   Value   `json:"value"`
		TypeURL *string `json:"type_url"`
	} `json:"parameter"`
	Type *string `json:"type"`
}

type Value struct {
	Amount       *int64  `json:"amount"`
	OwnerAddress *string `json:"owner_address"`
	ToAddress    *string `json:"to_address"`
}

type Meta struct {
	At          *int64  `json:"at"`
	Fingerprint *string `json:"fingerprint"`
	Links       struct {
		Next *string `json:"next"`
	} `json:"links"`
	PageSize *int `json:"page_size"`
}

type GetTxByAccountAddressFilter struct {
	OnlyConfirmed   *bool
	OnlyUnconfirmed *bool
	OnlyTo          *bool
	OnlyFrom        *bool
	Limit           *int64
	Fingerprint     *string
	OrderBy         *string
	MinTimestampMs  *int64
	MaxTimestampMs  *int64
}
