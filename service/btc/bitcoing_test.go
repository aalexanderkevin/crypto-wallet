package btc_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/service/btc"

	"github.com/stretchr/testify/require"
)

func TestServiceBtc_CheckAddress(t *testing.T) {
	t.Run("ShouldReturnTrue_WhenAddressIsValidOnTestnet", func(t *testing.T) {
		cfg := config.Instance()
		cfg.Bitcoin.Chain = "test3"
		btcSvc := btc.NewBitcoinImpl(cfg)

		// CODE UNDER TEST
		valid := btcSvc.CheckAddress(helper.Pointer("myAJasLvCqJJLkW2WzGr3S6Xkp4GKMTGPa"))

		// EXPECTATION
		require.True(t, valid)
	})

	t.Run("ShouldReturnFalse_WhenAddressIsNotValidOnTestnet", func(t *testing.T) {
		cfg := config.Instance()
		cfg.Bitcoin.Chain = "test3"
		btcSvc := btc.NewBitcoinImpl(cfg)

		// CODE UNDER TEST
		valid := btcSvc.CheckAddress(helper.Pointer("myAJasLvCqJJLkW2WzGr3S6Kkp4GKMTGPa"))

		// EXPECTATION
		require.False(t, valid)
	})

	t.Run("ShouldReturnTrue_WhenAddressIsValidOnMainNet", func(t *testing.T) {
		cfg := config.Instance()
		cfg.Bitcoin.Chain = "main"
		btcSvc := btc.NewBitcoinImpl(cfg)

		// CODE UNDER TEST
		valid := btcSvc.CheckAddress(helper.Pointer("bc1qrw5gztwp6vkycnfp25d3ejtdljr6g7che2qu7j"))

		// EXPECTATION
		require.True(t, valid)
	})

	t.Run("ShouldReturnFalse_WhenAddressIsNotValidOnMainNet", func(t *testing.T) {
		cfg := config.Instance()
		cfg.Bitcoin.Chain = "main"
		btcSvc := btc.NewBitcoinImpl(cfg)

		// CODE UNDER TEST
		valid := btcSvc.CheckAddress(helper.Pointer("bc1qrw5gztwp6vkycnkp25d3ejtdljr6g7che2qu7j"))

		// EXPECTATION
		require.False(t, valid)
	})

}

func TestServiceBtc_GetWallet(t *testing.T) {
	t.Run("ShouldGetBtcWallet", func(t *testing.T) {
		// INIT
		btcSvc := btc.NewBitcoinImpl(config.Instance())

		// seedPhrase := "yellow dolphin robot express road develop repair neutral rate tide economy section"
		seedPhrase := "east embrace bonus puzzle else have know fire essay unlock theme vibrant"

		// CODE UNDER TEST
		wallet, err := btcSvc.GetWallet(context.TODO(), &seedPhrase)
		require.NoError(t, err)

		// // EXPECTATION
		require.NotNil(t, wallet)
		require.NotNil(t, wallet.Address)
		require.Equal(t, "myAJasLvCqJJLkW2WzGr3S6Xkp4GKMTGPa", wallet.Address.EncodeAddress())
		require.Equal(t, "03eb2a6124a9994deb0602a68cf3868dca816be3235f4d1a94463473b92c72c5fa", fmt.Sprintf("%x", wallet.PublicKey.SerializeCompressed()))
	})

}

func TestServiceBtc_GetBalanceAddress(t *testing.T) {
	t.Run("ShouldGetBalance", func(t *testing.T) {
		btcSvc := btc.NewBitcoinImpl(config.Instance())

		// CODE UNDER TEST
		bal, err := btcSvc.GetBalance(context.TODO(), "myAJasLvCqJJLkW2WzGr3S6Xkp4GKMTGPa")

		// EXPECTATION
		require.NoError(t, err)
		require.Greater(t, bal.Int64(), int64(1))
	})
}

func TestServiceBtc_SendTx(t *testing.T) {
	t.Run("ShouldSendTx", func(t *testing.T) {
		// INIT
		btcSvc := btc.NewBitcoinImpl(config.Instance())

		// seedPhrase := "yellow dolphin robot express road develop repair neutral rate tide economy section"
		seedPhrase := "east embrace bonus puzzle else have know fire essay unlock theme vibrant"

		wallet, err := btcSvc.GetWallet(context.TODO(), &seedPhrase)
		require.NoError(t, err)

		// CODE UNDER TEST
		tx, err := btcSvc.SendTx(context.TODO(), wallet, &model.TxOpts{
			// To: helper.Pointer("myzJWXp5ywnjJiZRH6qoc6vioV9WQB9gJg"),
			To:     helper.Pointer("myAJasLvCqJJLkW2WzGr3S6Xkp4GKMTGPa"),
			Amount: big.NewInt(int64(2000)),
		})

		// EXPECTATION
		require.NoError(t, err)
		require.NotNil(t, tx)
	})

}

func TestServiceBtc_GetTx(t *testing.T) {
	t.Run("ShouldGetTx", func(t *testing.T) {
		// INIT
		btcSvc := btc.NewBitcoinImpl(config.Instance())

		// CODE UNDER TEST
		tx, err := btcSvc.GetTx(context.TODO(), "3c099630ee223a0f43bea658e4e4c88bfdea83f009a39ca05fdf764b4539e7f0")

		// EXPECTATION
		require.NoError(t, err)
		require.NotNil(t, tx)
	})

}

func TestServiceBtc_CreateWebhook(t *testing.T) {
	t.Run("ShouldReturnWebhookId", func(t *testing.T) {
		// INIT
		btcSvc := btc.NewBitcoinImpl(config.Instance())
		ctx := context.TODO()

		// CODE UNDER TEST
		webhook, err := btcSvc.CreateWebhookConfirmedTx(ctx, helper.Pointer("myzJWXp5ywnjJiZRH6qoc6vioV9WQB9gJg"))

		// EXPECTATION
		require.NoError(t, err)
		require.NotNil(t, webhook)
		require.NotEmpty(t, webhook.Address)

		err = btcSvc.DeleteWebhook(ctx, helper.Pointer(webhook.ID))
		require.NoError(t, err)
	})

}
