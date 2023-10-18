package usecase

import (
	"context"
	"fmt"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/container"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/repository"
	"github.com/aalexanderkevin/crypto-wallet/service"
)

type Wallet struct {
	config config.Config
	service.Bitcoin
	service.Ethereum
	service.Tron
	repository.Wallet
}

func NewWallet(c *container.Container) *Wallet {
	return &Wallet{
		config:   c.Config(),
		Bitcoin:  c.Bitcoin(),
		Ethereum: c.Ethereum(),
		Tron:     c.Tron(),
		Wallet:   c.WalletRepo(),
	}
}

func (w Wallet) CreateNewWallet(ctx context.Context, email *string) (wallet *model.Wallet, err error) {
	logger := helper.GetLogger(ctx).WithField("method", "Usecase.Wallet.CreateNewWallet")

	existingWallet, err := w.Wallet.Get(ctx, &repository.WalletGetFilter{
		Email: email,
	}, nil)
	if err == nil && existingWallet != nil {
		err = fmt.Errorf("wallet already exist")
		logger.WithError(err)
		return nil, err
	} else if err != nil && !model.IsNotFoundError(err) {
		logger.WithError(err).Warn("failed check existing wallet")
		return nil, err
	}

	seedPhrase := helper.GenerateSecureSeedPhrase()
	btcWallet, err := w.Bitcoin.GetWallet(ctx, &seedPhrase)
	if err != nil {
		logger.WithError(err).Warn("failed get new wallet by seedPhrase")
		return nil, err
	}

	ethWallet, err := w.Ethereum.GetWallet(ctx, &seedPhrase)
	if err != nil {
		logger.WithError(err).Warn("failed get new wallet by seedPhrase")
		return nil, err
	}

	trxWallet := w.Tron.GetWallet(ctx, &seedPhrase)

	wallet = &model.Wallet{}
	wallet.Email = email
	wallet.SeedPhrase = &seedPhrase
	wallet.BtcAddress = helper.Pointer(btcWallet.Address.EncodeAddress())
	wallet.EthAddress = helper.Pointer(ethWallet.Account.Address.Hex())
	wallet.EthAddress = trxWallet.Address

	wallet, err = w.Wallet.Add(ctx, wallet, &w.config.Service.SeedPhraseEncryptionKey)
	if err != nil {
		logger.WithError(err).Warn("failed insert wallet")
		return nil, err
	}

	return wallet, nil
}
