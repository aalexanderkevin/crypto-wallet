package usecase

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/container"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/repository"
	"github.com/aalexanderkevin/crypto-wallet/service"
	"github.com/ethereum/go-ethereum/common"
	"github.com/segmentio/ksuid"
)

type Transaction struct {
	config config.Config
	service.Bitcoin
	service.Ethereum
	service.Tron
	repository.Wallet

	btcTransactionRepo repository.Transaction
	ethTransactionRepo repository.Transaction
	trxTransactionRepo repository.Transaction

	sleepCheckPendingTrx      time.Duration
	sleepCheckConfirmationTrx time.Duration
}

func NewTransaction(c *container.Container) *Transaction {
	return &Transaction{
		config:             c.Config(),
		Bitcoin:            c.Bitcoin(),
		Ethereum:           c.Ethereum(),
		Tron:               c.Tron(),
		btcTransactionRepo: c.TransactionBtcRepo(),
		ethTransactionRepo: c.TransactionEthRepo(),
		trxTransactionRepo: c.TransactionTrxRepo(),
		Wallet:             c.WalletRepo(),

		sleepCheckPendingTrx:      5 * time.Second,
		sleepCheckConfirmationTrx: 1 * time.Minute,
	}
}

func (t Transaction) SendBitcoin(ctx context.Context, reqSend *model.SendToken) (txHash *string, err error) {
	logger := helper.GetLogger(ctx).WithField("method", "Usecase.Transaction.SendBitcoin")

	// validate receiver bitcoin address
	if valid := t.Bitcoin.CheckAddress(reqSend.ReceiverAddress); !valid {
		err = fmt.Errorf("invalid receiver bitcoin address")
		logger.WithError(err)
		return nil, err
	}

	// get the seedphrase of sender
	wallet, err := t.Wallet.Get(ctx, &repository.WalletGetFilter{
		Email: reqSend.Email,
	}, &t.config.Service.SeedPhraseEncryptionKey)
	if err != nil {
		logger.WithError(err).Warn("failed get wallet")
		return nil, err
	}

	// generate seedphrase to get btc wallet
	btcWallet, err := t.Bitcoin.GetWallet(ctx, wallet.SeedPhrase)
	if err != nil {
		logger.WithError(err).Warn("failed get btc wallet")
		return nil, err
	}

	// send token
	tx, err := t.Bitcoin.SendTx(ctx, btcWallet, &model.TxOpts{
		To:     reqSend.ReceiverAddress,
		Amount: big.NewInt(*reqSend.Amount),
	})
	if err != nil {
		logger.WithError(err).Warn("failed send btc")
		return nil, err
	}

	// upsert tx on database
	if _, err = t.btcTransactionRepo.Upsert(ctx, tx); err != nil {
		logger.WithError(err).Warn("Failed upsert")
	}

	return tx.Id, nil
}

func (t Transaction) SendTron(ctx context.Context, reqSend *model.SendToken) (txHash *string, err error) {
	logger := helper.GetLogger(ctx).WithField("method", "Usecase.Transaction.SendTron")

	// validate receiver tron address
	if valid := t.Tron.CheckAddress(*reqSend.ReceiverAddress); valid != nil {
		err = fmt.Errorf("invalid receiver tron address")
		logger.WithError(err)
		return nil, err
	}

	// get the seedphrase of sender
	wallet, err := t.Wallet.Get(ctx, &repository.WalletGetFilter{
		Email: reqSend.Email,
	}, &t.config.Service.SeedPhraseEncryptionKey)
	if err != nil {
		logger.WithError(err).Warn("failed get wallet")
		return nil, err
	}

	// generate seedphrase to get trx wallet
	trxWallet := t.Tron.GetWallet(ctx, wallet.SeedPhrase)

	// send token
	tx, err := t.Tron.SendTx(ctx, &model.TxOpts{
		To:     reqSend.ReceiverAddress,
		Amount: big.NewInt(*reqSend.Amount),
	}, trxWallet)
	if err != nil {
		logger.WithError(err).Warn("failed send trx")
		return nil, err
	}

	transactionHash := common.BytesToHash(tx.Txid)

	transaction := &model.Transaction{
		Id:              helper.Pointer(transactionHash.Hex()[2:]),
		SenderAddress:   []string{*wallet.TrxAddress},
		ReceiverAddress: []string{*reqSend.ReceiverAddress},
		Amount:          reqSend.Amount,
		Status:          helper.Pointer("pending"),
	}

	// open new thread to check transaction success
	go t.checkTransactionTrx(ctx, transaction)

	return transaction.Id, nil
}

func (t Transaction) checkTransactionTrx(ctx context.Context, transaction *model.Transaction) {
	reqId := ctx.Value(helper.ContextKeyRequestId)
	if reqId == nil {
		reqId = ctx.Value(string(helper.ContextKeyRequestId))
		if reqId == nil {
			reqId = ksuid.New().String()
		}
	}
	ctx = helper.ContextWithRequestId(context.Background(), reqId.(string))
	logger := helper.GetLogger(ctx).WithField("method", "Usecase.Wallet.checkTransaction")

	for {
		txInfo, err := t.Tron.GetTx(ctx, *transaction.Id)
		if err != nil {
			if err.Error() != "transaction info not found" {
				logger.WithError(err).Warn("failed to get trx")
				return
			}
		}

		if txInfo != nil {
			transaction.Fee = &txInfo.Fee
			transaction.ReceivedAt = helper.Pointer(time.UnixMilli(txInfo.BlockTimeStamp))
			transaction.Block = &txInfo.BlockNumber

			// get the current block
			currentBlock, err := t.Tron.GetCurrentBlock(ctx)
			if err != nil {
				logger.WithError(err).Warn("Failed get current block trx")
				return
			}

			// calucate the confirmation block
			if currentBlock != nil && transaction.Block != nil {
				transaction.Confirmation = helper.Pointer(*currentBlock - *transaction.Block)
			}

			// upsert tx on database
			if _, err = t.trxTransactionRepo.Upsert(ctx, transaction); err != nil {
				logger.WithError(err).Warn("Failed upsert transaction trx")
				return
			}

			break
		}

		time.Sleep(t.sleepCheckPendingTrx)
	}

	for {
		unconfirmed, err := t.Tron.GetUnconfirmedTxAddress(ctx, &transaction.SenderAddress[0])
		if err != nil {
			logger.WithError(err).Warn("failed to GetUnconfirmedTxAddress")
			return
		}

		for _, tx := range unconfirmed.Data {
			if *tx.TxID == *transaction.Id {
				time.Sleep(t.sleepCheckPendingTrx)
				continue
			}
		}

		confirmed, err := t.Tron.GetConfirmedTxAddress(ctx, &transaction.SenderAddress[0])
		if err != nil {
			logger.WithError(err).Warn("failed to GetConfirmedTxAddress")
			return
		}

		for _, tx := range confirmed.Data {
			if *tx.TxID == *transaction.Id {
				// get the current block
				currentBlock, err := t.Tron.GetCurrentBlock(ctx)
				if err != nil {
					logger.WithError(err).Warn("Failed get current block trx")
					return
				}

				// calucate the confirmation block
				if currentBlock != nil && transaction.Block != nil {
					transaction.Confirmation = helper.Pointer(*currentBlock - *transaction.Block)
				}

				transaction.Status = helper.Pointer("success")

				// upsert tx on database
				if _, err = t.trxTransactionRepo.Upsert(ctx, transaction); err != nil {
					logger.WithError(err).Warn("Failed upsert transaction trx")
				}
				return
			}
		}

		time.Sleep(t.sleepCheckConfirmationTrx)
	}
}

func (t Transaction) SendEthereum(ctx context.Context, reqSend *model.SendToken) (txHash *string, err error) {
	logger := helper.GetLogger(ctx).WithField("method", "Usecase.Transaction.SendEthereum")

	// validate receiver ethereum address
	if err := t.Ethereum.CheckAddress(*reqSend.ReceiverAddress); err != nil {
		err = fmt.Errorf("invalid receiver ethereum address")
		logger.WithError(err)
		return nil, err
	}

	// get the seedphrase of sender
	wallet, err := t.Wallet.Get(ctx, &repository.WalletGetFilter{
		Email: reqSend.Email,
	}, &t.config.Service.SeedPhraseEncryptionKey)
	if err != nil {
		logger.WithError(err).Warn("failed get wallet")
		return nil, err
	}

	// generate seedphrase to get trx wallet
	ethWallet, err := t.Ethereum.GetWallet(ctx, wallet.SeedPhrase)
	if err != nil {
		logger.WithError(err).Warn("failed get eth wallet")
		return nil, err
	}

	// send token
	tx, err := t.Ethereum.SendTx(ctx, &model.TxOpts{
		To:     reqSend.ReceiverAddress,
		Amount: big.NewInt(*reqSend.Amount),
	}, ethWallet)
	if err != nil {
		logger.WithError(err).Warn("failed send trx")
		return nil, err
	}

	transaction := &model.Transaction{
		Id:              helper.Pointer(tx.Hash().Hex()),
		SenderAddress:   []string{*wallet.EthAddress},
		ReceiverAddress: []string{*reqSend.ReceiverAddress},
		Amount:          reqSend.Amount,
		Confirmation:    helper.Pointer[int64](0),
		Fee:             helper.Pointer(tx.GasPrice().Int64() * int64(tx.Gas())),
		Status:          helper.Pointer("pending"),
	}

	// open new thread to check transaction success
	go t.CheckTransactionEth(ctx, transaction)

	_, err = t.ethTransactionRepo.Upsert(ctx, transaction)
	if err != nil {
		logger.WithError(err).Warn("failed Upsert eth transaction")
		return nil, err
	}

	return transaction.Id, nil
}

func (t Transaction) CheckTransactionEth(ctx context.Context, transaction *model.Transaction) {
	reqId := ctx.Value(helper.ContextKeyRequestId)
	if reqId == nil {
		reqId = ctx.Value(string(helper.ContextKeyRequestId))
		if reqId == nil {
			reqId = ksuid.New().String()
		}
	}
	ctx = helper.ContextWithRequestId(context.Background(), reqId.(string))
	logger := helper.GetLogger(ctx).WithField("method", "Usecase.Wallet.checkTransactionEth")

	txHash := helper.Pointer(common.HexToHash(*transaction.Id))
	tx := &model.Transaction{Id: transaction.Id}
	for {

		var isPending *bool
		var err error
		tx, isPending, err = t.Ethereum.GetTransactionPending(ctx, txHash)
		if err != nil {
			logger.WithError(err).Warn("failed to GetTransactionPending")
			return
		}

		if isPending != nil && !*isPending {
			break
		}

		time.Sleep(t.sleepCheckPendingTrx)
	}

	for {
		blockDetails, err := t.Ethereum.GetBlockInformation(ctx, txHash)
		if err != nil {
			logger.WithError(err).Warn("failed to GetUnconfirmedTxAddress")
			return
		}

		tx.ReceivedAt = blockDetails.ReceivedAt
		tx.Block = blockDetails.Block
		if blockDetails.Status != nil {
			tx.Status = blockDetails.Status
		}
		if *blockDetails.Confirmation > *tx.Confirmation {
			tx.Confirmation = blockDetails.Confirmation

			_, err = t.ethTransactionRepo.Upsert(ctx, tx)
			if err != nil {
				logger.WithError(err).Warn("failed Upsert eth transaction")
				return
			}
		}

		if *tx.Status == "success" {
			break
		}

		time.Sleep(t.sleepCheckConfirmationTrx)
	}
}
