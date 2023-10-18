package repository

import (
	"context"

	"github.com/aalexanderkevin/crypto-wallet/model"
)

type Wallet interface {
	Add(ctx context.Context, wallet *model.Wallet, encryptionKey *string) (*model.Wallet, error)
	Get(ctx context.Context, filter *WalletGetFilter, encryptionKey *string) (*model.Wallet, error)
	Update(ctx context.Context, id string, wallet *model.Wallet) (*model.Wallet, error)
}

type WalletGetFilter struct {
	Id    *string
	Email *string
}
