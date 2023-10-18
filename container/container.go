package container

import (
	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/repository"
	"github.com/aalexanderkevin/crypto-wallet/service"

	"gorm.io/gorm"
)

type Container struct {
	config config.Config
	db     *gorm.DB

	//svc
	ethereum service.Ethereum
	bitcoin  service.Bitcoin
	tron     service.Tron
	redis    service.Cache

	// repo
	walletRepo         repository.Wallet
	transactionBtcRepo repository.Transaction
	transactionEthRepo repository.Transaction
	transactionTrxRepo repository.Transaction
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) Config() config.Config {
	return c.config
}

func (c *Container) SetConfig(config config.Config) {
	c.config = config
}

func (c *Container) Db() *gorm.DB {
	return c.db
}

func (c *Container) SetDb(db *gorm.DB) {
	c.db = db
}

func (c *Container) Redis() service.Cache {
	return c.redis
}

func (c *Container) SetRedis(redis service.Cache) {
	c.redis = redis
}

func (c *Container) Ethereum() service.Ethereum {
	return c.ethereum
}

func (c *Container) SetEthereum(ethereum service.Ethereum) {
	c.ethereum = ethereum
}

func (c *Container) Bitcoin() service.Bitcoin {
	return c.bitcoin
}

func (c *Container) SetBitcoin(bitcoin service.Bitcoin) {
	c.bitcoin = bitcoin
}

func (c *Container) Tron() service.Tron {
	return c.tron
}

func (c *Container) SetTron(tron service.Tron) {
	c.tron = tron
}

func (c *Container) WalletRepo() repository.Wallet {
	return c.walletRepo
}

func (c *Container) SetWalletRepo(walletRepo repository.Wallet) {
	c.walletRepo = walletRepo
}

func (c *Container) TransactionEthRepo() repository.Transaction {
	return c.transactionEthRepo
}

func (c *Container) SetTransactionEthRepo(transactionEthRepo repository.Transaction) {
	c.transactionEthRepo = transactionEthRepo
}

func (c *Container) TransactionBtcRepo() repository.Transaction {
	return c.transactionBtcRepo
}

func (c *Container) SetTransactionBtcRepo(transactionBtcRepo repository.Transaction) {
	c.transactionBtcRepo = transactionBtcRepo
}

func (c *Container) TransactionTrxRepo() repository.Transaction {
	return c.transactionTrxRepo
}

func (c *Container) SetTransactionTrxRepo(transactionTrxRepo repository.Transaction) {
	c.transactionTrxRepo = transactionTrxRepo
}
