package eth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"regexp"
	"time"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/service"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

type EthereumImpl struct {
	client *ethclient.Client
	ws     *ethclient.Client
	config config.Ethereum
}

func NewEthereumImpl(config config.Config) service.Ethereum {
	client, err := ethclient.Dial(config.Ethereum.NetUrl)

	if err != nil {
		panic(fmt.Sprintf("error connect to eth client: %s, with error %v", config.Ethereum.NetUrl, err))
	}

	ws, err := ethclient.Dial("wss://sepolia.infura.io/ws/v3/282be59eb719440b89fe8168d85003fb")
	if err != nil {
		panic(fmt.Sprintf("error connect to eth client: %s, with error %v", config.Ethereum.NetUrl, err))
	}

	return &EthereumImpl{
		client: client,
		config: config.Ethereum,
		ws:     ws,
	}
}

func (e *EthereumImpl) Close() {
	e.client.Close()
	e.ws.Close()
}

func (e *EthereumImpl) GetWallet(ctx context.Context, seedPhrase *string) (*model.EthHdWallet, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Ethereum.GetWallet")

	ethWallet, err := hdwallet.NewFromMnemonic(*seedPhrase)
	if err != nil {
		logger.WithError(err).Warn("Failed create ethwallet from mnemonic")
		return nil, err
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := ethWallet.Derive(path, false)
	if err != nil {
		logger.WithError(err).Warn("Failed derive eth path")
		return nil, err
	}

	wallet := &model.EthHdWallet{
		Wallet:  ethWallet,
		Account: &account,
	}
	return wallet, nil
}

func (e *EthereumImpl) SendTx(ctx context.Context, txOpts *model.TxOpts, wallet *model.EthHdWallet) (tx *types.Transaction, err error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Ethereum.SendTx")

	if err := e.CheckAddress(*txOpts.To); err != nil {
		logger.WithError(err).Warn("Failed validate eth address")
		return nil, err
	}

	tx, err = e.createTx(ctx, txOpts, wallet)
	if err != nil {
		logger.WithError(err).Warn("Failed createTx")
		return nil, err
	}

	err = e.client.SendTransaction(ctx, tx)
	if err != nil {
		logger.WithError(err).Warn("Failed sendTransaction ethereum")
		return
	}

	return
}

func (e *EthereumImpl) GetBalance(ctx context.Context, fromAddress common.Address) (balance *big.Int, err error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Ethereum.GetBalance")

	balance, err = e.client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		logger.WithError(err).Warn("Failed balanceAt Ethereum")
		return nil, err
	}

	return
}

func (e *EthereumImpl) getNonce(ctx context.Context, fromAddress common.Address) (nonce uint64, err error) {
	return e.client.PendingNonceAt(ctx, fromAddress)
}

func (e *EthereumImpl) getGasLimit(ctx context.Context, fromAddress common.Address, txOpts model.TxOpts) (gasLimit uint64, err error) {
	toAddress := common.HexToAddress(*txOpts.To)

	return e.client.EstimateGas(ctx, ethereum.CallMsg{
		From:  fromAddress,
		To:    helper.Pointer(toAddress),
		Value: txOpts.Amount,
	})
}

func (e *EthereumImpl) getGasPrice(ctx context.Context) (gasPrice *big.Int, err error) {
	return e.client.SuggestGasPrice(ctx)
}

func (e *EthereumImpl) createTx(ctx context.Context, txOpts *model.TxOpts, wallet *model.EthHdWallet) (*types.Transaction, error) {
	nonce, err := e.getNonce(ctx, wallet.Account.Address)
	if err != nil {
		return nil, err
	}

	gasLimit, err := e.getGasLimit(ctx, wallet.Account.Address, *txOpts)
	if err != nil {
		return nil, err
	}

	gasPrice, err := e.getGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	balance, err := e.GetBalance(ctx, wallet.Account.Address)
	if err != nil {
		return nil, err
	}

	ok := checkValueEnough(txOpts.Amount, gasPrice, gasLimit, balance)
	if !ok {
		return nil, fmt.Errorf("error not enough balance")
	}

	toAddress := common.HexToAddress(*txOpts.To)
	baseTx := &types.DynamicFeeTx{
		Nonce:     uint64(nonce),
		GasTipCap: gasPrice,
		GasFeeCap: gasPrice,
		Gas:       uint64(gasLimit),
		To:        helper.Pointer(toAddress),
		Value:     txOpts.Amount,
		Data:      nil,
	}

	chainID, err := e.client.NetworkID(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	privateKey, err := wallet.Wallet.PrivateKey(*wallet.Account)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	signedTx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(chainID), baseTx)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return signedTx, nil

}

func checkValueEnough(value *big.Int, gasPrice *big.Int, gasLimit uint64, balance *big.Int) bool {
	tvalue := big.NewInt(0).Set(value)
	tgasPrice := big.NewInt(0).Set(gasPrice)
	return tvalue.Add(tvalue, tgasPrice.Mul(tgasPrice, big.NewInt(int64(gasLimit)))).Cmp(balance) != 1
}

func (e *EthereumImpl) GetTx(ctx context.Context, txHash *common.Hash) (*model.Transaction, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Ethereum.GetTx")

	res, isPending, err := e.GetTransactionPending(ctx, txHash)
	if err != nil {
		logger.WithError(err).Warn("Failed GetTransactionPending")
		return nil, err
	}

	// return if transaction still pending
	if isPending != nil && *isPending {
		return res, nil
	}

	blockInformation, err := e.GetBlockInformation(ctx, txHash)
	if err != nil {
		logger.WithError(err).Warn("Failed GetBlockInformation")
		return nil, err
	}

	res.Block = blockInformation.Block
	res.ReceivedAt = blockInformation.ReceivedAt
	res.Confirmation = blockInformation.Confirmation
	res.Status = blockInformation.Status

	return res, nil
}

type EtherscanTransaction struct {
	BlockNumber string `json:"blockNumber"`
	TimeStamp   string `json:"timeStamp"`
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	GasPrice    string `json:"gasPrice"`
	Gas         string `json:"gas"`
}

func (e *EthereumImpl) GetCurrentBlock(ctx context.Context) (*int64, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Ethereum.GetCurrentBlock")

	blockNum, err := e.client.BlockNumber(ctx)
	if err != nil {
		logger.WithError(err).Warn("Failed get BlockNumber")
		return nil, err
	}

	return helper.Pointer(int64(blockNum)), nil
}

func (e *EthereumImpl) CheckAddress(address string) error {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(address) {
		return errors.New("invalid address")
	}

	return nil
}

func (e *EthereumImpl) GetBlockInformation(ctx context.Context, txHash *common.Hash) (*model.Transaction, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Ethereum.GetTx")
	res := &model.Transaction{}

	receipt, err := e.client.TransactionReceipt(context.Background(), *txHash)
	if err != nil {
		logger.WithError(err).Warn("Failed get TransactionReceipt")
		return nil, err
	}

	blck, err := e.client.BlockByHash(ctx, receipt.BlockHash)
	if err != nil {
		logger.WithError(err).Warn("Failed get BlockByHash")
		return nil, err
	}

	currentBlock, err := e.GetCurrentBlock(ctx)
	if err != nil {
		logger.WithError(err).Warn("Failed GetCurrentBlock")
		return nil, err
	}

	res.Block = helper.Pointer(receipt.BlockNumber.Int64())
	res.ReceivedAt = helper.Pointer(time.Unix(int64(blck.Time()), 0))
	res.Confirmation = helper.Pointer(int64(*currentBlock) - *res.Block)

	if *res.Confirmation > 12 {
		res.Status = helper.Pointer("success")
	}

	return res, nil
}

func (e *EthereumImpl) GetTransactionPending(ctx context.Context, txHash *common.Hash) (*model.Transaction, *bool, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Ethereum.GetTransactionPending")

	res := &model.Transaction{Id: helper.Pointer(txHash.Hex())}
	res.Status = helper.Pointer("pending")
	res.Confirmation = helper.Pointer[int64](0)

	tx, isPending, err := e.client.TransactionByHash(ctx, *txHash)
	if err != nil {
		logger.WithError(err).Warn("Failed get TransactionByHash")
		return nil, nil, err
	}

	res.ReceiverAddress = []string{tx.To().Hex()}
	res.Amount = helper.Pointer(tx.Value().Int64())
	res.Fee = helper.Pointer(tx.GasPrice().Int64() * int64(tx.Gas()))

	chainId, err := e.client.ChainID(ctx)
	if err != nil {
		logger.WithError(err).Warn("Failed get ChainID")
		return nil, nil, err
	}

	// get sender address
	from, err := types.Sender(types.LatestSignerForChainID(chainId), tx)
	if err != nil {
		logger.WithError(err).Warn("Failed get transaction sender")
		return nil, nil, err
	}
	res.SenderAddress = []string{from.Hex()}

	if isPending {
		return res, helper.Pointer(true), nil
	}

	return res, helper.Pointer(false), nil
}

func (e *EthereumImpl) SubscribePendingTransactions(ctx context.Context) (subs *rpc.ClientSubscription, txch chan *types.Transaction, err error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Ethereum.RunningWatcher")
	gcli := gethclient.New(e.ws.Client())

	txch = make(chan *types.Transaction, 100)
	subs, err = gcli.SubscribeFullPendingTransactions(context.Background(), txch)
	if err != nil {
		logger.WithError(err).Warn("Failed to SubscribeFullPendingTransactions")
	}

	return
}
