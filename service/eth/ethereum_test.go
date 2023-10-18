//go:build integration
// +build integration

package eth_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/service/eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestServiceEth_GetWallet(t *testing.T) {
	t.Run("ShouldGetEthWallet", func(t *testing.T) {
		// INIT
		ethSvc := eth.NewEthereumImpl(config.Instance())
		defer ethSvc.Close()

		seedPhrase := "yellow dolphin robot express road develop repair neutral rate tide economy section"

		// CODE UNDER TEST
		wallet, err := ethSvc.GetWallet(context.TODO(), &seedPhrase)
		require.NoError(t, err)

		// EXPECTATION
		require.NotNil(t, wallet)
		require.NotNil(t, wallet.Account)
		require.Equal(t, "0x8b0FD6E2530CAA4f4cd0BdaA5c05d06AFd0D538f", wallet.Account.Address.Hex())
	})

}

func TestServiceEth_GetBalance(t *testing.T) {
	t.Run("ShouldGetBalance", func(t *testing.T) {
		// INIT
		ethSvc := eth.NewEthereumImpl(config.Instance())
		defer ethSvc.Close()

		seedPhrase := "yellow dolphin robot express road develop repair neutral rate tide economy section"

		wallet, err := ethSvc.GetWallet(context.TODO(), &seedPhrase)
		require.NoError(t, err)
		require.NotNil(t, wallet)
		require.NotNil(t, wallet.Account)
		require.Equal(t, "0x8b0FD6E2530CAA4f4cd0BdaA5c05d06AFd0D538f", wallet.Account.Address.Hex())

		// CODE UNDER TEST
		balance, err := ethSvc.GetBalance(context.TODO(), wallet.Account.Address)
		require.NoError(t, err)

		// EXPECTATION
		require.NoError(t, err)
		require.NotNil(t, balance)
		require.Greater(t, 0.5, *balance)
	})
}

func TestServiceEth_SendTx(t *testing.T) {
	t.Run("ShouldSendTx", func(t *testing.T) {
		// INIT
		ethSvc := eth.NewEthereumImpl(config.Instance())
		defer ethSvc.Close()

		seedPhrase := "yellow dolphin robot express road develop repair neutral rate tide economy section"
		// seedPhrase := "east embrace bonus puzzle else have know fire essay unlock theme vibrant"

		wallet, err := ethSvc.GetWallet(context.TODO(), &seedPhrase)
		require.NoError(t, err)
		require.NotNil(t, wallet)
		require.NotNil(t, wallet.Account)
		require.Equal(t, "0x8b0FD6E2530CAA4f4cd0BdaA5c05d06AFd0D538f", wallet.Account.Address.Hex())

		to := helper.HexToAddress("0x12015F59614e650a1554b8CFF444580207c7a5c6")
		value := helper.EthToWei(0.01111)

		// CODE UNDER TEST
		tx, err := ethSvc.SendTx(context.TODO(), &model.TxOpts{
			To:     helper.Pointer(to.Hex()),
			Amount: value,
		}, wallet)

		fmt.Println(tx.Hash())
		fmt.Println(tx.Value())
		fmt.Println(tx.GasPrice())

		fmt.Println(tx.Time())

		// EXPECTATION
		require.NoError(t, err)
		require.NotNil(t, tx)

		txa, err := ethSvc.GetTx(context.TODO(), helper.Pointer(tx.Hash()))
		require.NoError(t, err)
		require.NotNil(t, txa)

	})
}

func TestServiceEth_GetTx(t *testing.T) {
	t.Run("ShouldGetTx", func(t *testing.T) {
		// INIT
		ethSvc := eth.NewEthereumImpl(config.Instance())
		defer ethSvc.Close()

		txHash := common.HexToHash("0x47379d2fef6de3fbfad1fc023391c8b35ec15639dd1bcb97f2709392505bb3ba")
		// CODE UNDER TEST
		tx, err := ethSvc.GetTx(context.TODO(), &txHash)

		// EXPECTATION
		require.NoError(t, err)
		require.NotNil(t, tx)
		require.Equal(t, "success", *tx.Status)
		require.Equal(t, int64(111000000000000), *tx.Amount)
		require.Equal(t, int64(222702382203000), *tx.Fee)
		require.Equal(t, "0x8b0FD6E2530CAA4f4cd0BdaA5c05d06AFd0D538f", tx.SenderAddress[0])
		require.Equal(t, "0x12015F59614e650a1554b8CFF444580207c7a5c6", tx.ReceiverAddress[0])
		require.Greater(t, *tx.Confirmation, int64(12))
	})
}
