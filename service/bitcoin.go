package service

import (
	"context"
	"math/big"

	"github.com/aalexanderkevin/crypto-wallet/model"

	"github.com/blockcypher/gobcy/v2"
)

type Bitcoin interface {
	CheckAddress(address *string) bool
	GetWallet(ctx context.Context, seedPhrase *string) (*model.BtcHdWallet, error)
	GetBalance(ctx context.Context, address string) (*big.Int, error)
	SendTx(ctx context.Context, wallet *model.BtcHdWallet, txOpts *model.TxOpts) (*model.Transaction, error)
	GetTx(ctx context.Context, txhash string) (*gobcy.TX, error)
	CreateWebhookConfirmedTx(ctx context.Context, address *string) (*gobcy.Hook, error)
	DeleteWebhook(ctx context.Context, id *string) error
}
