package model

import (
	"crypto/ecdsa"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

type Wallet struct {
	Id         *string
	Email      *string
	SeedPhrase *string
	BtcAddress *string
	EthAddress *string
	TrxAddress *string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

type EthHdWallet struct {
	Wallet  *hdwallet.Wallet
	Account *accounts.Account
}

type TrxHdWallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    *string
}

type BtcHdWallet struct {
	PublicKey  *btcec.PublicKey
	Address    *btcutil.AddressPubKey
	PrivateKey *btcec.PrivateKey
	Wif        *btcutil.WIF
}
