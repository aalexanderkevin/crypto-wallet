package repository

import (
	"context"

	"github.com/aalexanderkevin/crypto-wallet/model"
)

type Transaction interface {
	Add(ctx context.Context, trx *model.Transaction) (*model.Transaction, error)
	Update(ctx context.Context, id string, trx *model.Transaction) (*model.Transaction, error)
	Upsert(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
	Get(ctx context.Context, filter *TransactionGetFilter) (*model.Transaction, error)
}

type TransactionGetFilter struct {
	Id              *string
	SenderAddress   *string
	ReceiverAddress *string
	Status          *string
}
