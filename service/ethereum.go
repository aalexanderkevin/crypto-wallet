package service

import (
	"context"
	"math/big"

	"github.com/aalexanderkevin/crypto-wallet/model"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type Ethereum interface {
	Close()
	GetWallet(ctx context.Context, seedPhrase *string) (*model.EthHdWallet, error)
	GetBalance(ctx context.Context, fromAddress common.Address) (*big.Int, error)
	SendTx(ctx context.Context, txOpts *model.TxOpts, wallet *model.EthHdWallet) (*types.Transaction, error)
	GetTx(ctx context.Context, txHash *common.Hash) (*model.Transaction, error)
	GetCurrentBlock(ctx context.Context) (*int64, error)
	CheckAddress(address string) error
	GetBlockInformation(ctx context.Context, txHash *common.Hash) (*model.Transaction, error)
	GetTransactionPending(ctx context.Context, txHash *common.Hash) (*model.Transaction, *bool, error)
	SubscribePendingTransactions(ctx context.Context) (subs *rpc.ClientSubscription, txch chan *types.Transaction, err error)
}
