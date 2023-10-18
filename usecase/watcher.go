package usecase

import (
	"context"
	"time"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/container"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/repository"
	"github.com/aalexanderkevin/crypto-wallet/service"
	"golang.org/x/exp/slices"
)

type Watcher struct {
	config config.Config
	service.Bitcoin
	service.Ethereum
	service.Tron
	service.Cache

	btcTransactionRepo repository.Transaction
	ethTransactionRepo repository.Transaction
	trxTransactionRepo repository.Transaction
	repository.Wallet

	usecaseTransaction Transaction
}

func NewWatcher(c *container.Container, t Transaction) *Watcher {
	return &Watcher{
		config:             c.Config(),
		Bitcoin:            c.Bitcoin(),
		Ethereum:           c.Ethereum(),
		Tron:               c.Tron(),
		Wallet:             c.WalletRepo(),
		btcTransactionRepo: c.TransactionBtcRepo(),
		ethTransactionRepo: c.TransactionEthRepo(),
		trxTransactionRepo: c.TransactionTrxRepo(),
		Cache:              c.Redis(),
		usecaseTransaction: t,
	}
}

func (w *Watcher) TriggerWatcherEth(ctx context.Context, email *string) (*string, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Usecase.Watcher.TriggerWatcherEth")

	// get the seedphrase of sender
	wallet, err := w.Wallet.Get(ctx, &repository.WalletGetFilter{
		Email: email,
	}, &w.config.Service.SeedPhraseEncryptionKey)
	if err != nil {
		logger.WithError(err).Warn("failed get wallet")
		return nil, err
	}

	res, err := w.Cache.SetList(ctx, "address", *wallet.EthAddress, helper.Pointer(2*time.Minute))
	if err != nil {
		logger.WithError(err).Warn("failed set address on cache")
		return nil, err
	}

	if res == 1 {
		go w.RunningWatcherEth(context.Background(), helper.Pointer("address"))
	}

	return wallet.EthAddress, nil
}

func (w *Watcher) RunningWatcherEth(ctx context.Context, cacheKeyList *string) {
	logger := helper.GetLogger(ctx).WithField("method", "Usecase.Watcher.RunningWatcherEth")

subscribe:
	subs, txch, err := w.SubscribePendingTransactions(ctx)
	if err != nil {
		logger.WithError(err).Warn("Failed to SubscribeFullPendingTransactions")
		return
	}

	defer subs.Unsubscribe()

	for {
		select {
		case tx := <-txch:
			addresses, err := w.Cache.GetList(ctx, *cacheKeyList)
			if err != nil {
				logger.WithError(err).Warn("Failed to get list address from cache")
				return
			}

			if tx.To() != nil && slices.Contains(addresses, tx.To().Hex()) {
				go w.usecaseTransaction.CheckTransactionEth(ctx, &model.Transaction{
					Id:              helper.Pointer(tx.Hash().Hex()),
					ReceiverAddress: []string{tx.To().Hex()},
				})
			}

		case err := <-subs.Err():
			logger.WithError(err).Warn("subscribe client connection is closed unexpectedly")
			goto subscribe
		}
	}
}

func (w *Watcher) TriggerWatcherTrx(ctx context.Context, email *string) (*string, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Usecase.Watcher.TriggerWatcherTrx")

	// get the seedphrase of sender
	wallet, err := w.Wallet.Get(ctx, &repository.WalletGetFilter{
		Email: email,
	}, &w.config.Service.SeedPhraseEncryptionKey)
	if err != nil {
		logger.WithError(err).Warn("failed get wallet")
		return nil, err
	}

	go w.RunningWatcherTrx(context.Background(), wallet.TrxAddress)

	return wallet.TrxAddress, nil
}

func (w *Watcher) RunningWatcherTrx(ctx context.Context, trxAddress *string) {
	logger := helper.GetLogger(ctx).WithField("method", "Usecase.Watcher.RunningWatcherTrx")
	newctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	startWatcher := time.Now().UnixMilli()

subscribe:
	for {
		select {
		case <-newctx.Done():
			cancel()
			return
		default:
			var fingerprint *string

			// get confirmed transaction
			trx, err := w.Tron.GetTxByAccountAddress(ctx, trxAddress, &service.GetTxByAccountAddressFilter{
				OnlyConfirmed:  helper.Pointer(true),
				OnlyTo:         helper.Pointer(true),
				OrderBy:        helper.Pointer("block_timestamp,asc"),
				MinTimestampMs: helper.Pointer(startWatcher),
				Fingerprint:    fingerprint,
			})
			if err != nil {
				logger.WithError(err).Warn("failed get confirmed transaction")
				goto subscribe
			}

			for _, data := range trx.Data {
				var amount *int64
				var to *string
				var from *string
				if data.RawData != nil && len(data.RawData.Contract) > 0 {
					amount = data.RawData.Contract[0].Parameter.Value.Amount
					to = helper.ToTrxAddress(*data.RawData.Contract[0].Parameter.Value.ToAddress)
					from = helper.ToTrxAddress(*data.RawData.Contract[0].Parameter.Value.OwnerAddress)
				}

				_, err := w.trxTransactionRepo.Upsert(ctx, &model.Transaction{
					Id:              data.TxID,
					SenderAddress:   []string{*from},
					ReceiverAddress: []string{*to},
					Amount:          amount,
					ReceivedAt:      helper.Pointer(time.UnixMilli(*data.BlockTimestamp)),
					Fee:             data.NetFee,
					Block:           data.BlockNumber,
					Status:          helper.Pointer("success"),
				})
				if err != nil {
					logger.WithError(err).Warn("failed upsert trx transaction")
					goto subscribe
				}

				startWatcher = *data.BlockTimestamp + 1
			}

			if trx.Meta != nil && trx.Meta.Fingerprint != nil {
				fingerprint = trx.Meta.Fingerprint
				goto subscribe
			}

			time.Sleep(1 * time.Minute)
		}
	}

}
