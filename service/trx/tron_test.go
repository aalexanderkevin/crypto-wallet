//go:build integration
// +build integration

package trx_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/service/trx"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/mr-tron/base58"

	"github.com/stretchr/testify/require"
)

func s256(s []byte) []byte {
	h := sha256.New()
	h.Write(s)
	bs := h.Sum(nil)
	return bs
}

func TestServiceTron_GetWallet(t *testing.T) {
	t.Run("ShouldGetWallet", func(t *testing.T) {

		addb, _ := hex.DecodeString("418c11ef4f7006a1c885cd62209f0b996b51030049")
		hash11 := s256(s256(addb))
		secrets := hash11[:4]
		addb = append(addb, secrets...)

		aa := helper.Pointer(base58.Encode([]byte(addb)))
		fmt.Println(aa)
		// INIT
		tronSvc := trx.NewTronImpl(config.Instance())
		defer tronSvc.Close()

		seedPhrase := "yellow dolphin robot express road develop repair neutral rate tide economy section"

		// CODE UNDER TEST
		trxWallet := tronSvc.GetWallet(context.Background(), &seedPhrase)

		// EXPECTATION
		require.NotNil(t, trxWallet)
		require.Equal(t, "04381b7ee8dc8d159f5ee80abf61702aed2644757ca81e298303a6813bfa36ecf60bd029f9cb2d224eb6405b36af182dccd1d9423a404af4a0cbe9a0518ae27a5d", hexutil.Encode(crypto.FromECDSAPub(trxWallet.PublicKey))[2:])
		require.Equal(t, "d73c4a10cf55567a18d671d52ce82599681787875326268ceca81a3ab57788eb", hexutil.Encode(crypto.FromECDSA(trxWallet.PrivateKey))[2:])
		require.Equal(t, "TNjq63hm9JfqQYRRwVAtS84PRy1Ty6CU5U", *trxWallet.Address)
	})

}

func TestServiceTrx_GetBalance(t *testing.T) {
	t.Run("ShouldGetBalance", func(t *testing.T) {
		// INIT
		trxSvc := trx.NewTronImpl(config.Instance())
		defer trxSvc.Close()

		// CODE UNDER TEST
		balance, err := trxSvc.GetBalance(context.TODO(), helper.Pointer("TNjq63hm9JfqQYRRwVAtS84PRy1Ty6CU5U"))
		// EXPECTATION
		require.NoError(t, err)
		require.NotNil(t, balance)
		require.Greater(t, *balance, int64(1))
	})
}

func TestServiceTrx_SendTx(t *testing.T) {
	t.Run("ShouldSendTx", func(t *testing.T) {
		// INIT
		tronSvc := trx.NewTronImpl(config.Instance())
		defer tronSvc.Close()

		seedPhrase := "yellow dolphin robot express road develop repair neutral rate tide economy section"
		trxWallet := tronSvc.GetWallet(context.Background(), &seedPhrase)

		// CODE UNDER TEST
		tx, err := tronSvc.SendTx(context.TODO(), &model.TxOpts{
			To:     helper.Pointer("TKTekM9b4m3ZtwfynqNj8TDgYpfQnEA2mr"),
			Amount: big.NewInt(100000),
		}, trxWallet)

		// EXPECTATION
		require.NoError(t, err)
		require.NotNil(t, tx)

		fmt.Println("timestamp tx : ")
		fmt.Println(tx.Transaction.RawData.Timestamp)
		hash := common.BytesToHash(tx.Txid)

		time.Sleep(3 * time.Second)
		txInfo, err := tronSvc.GetTx(context.TODO(), hash.Hex())
		require.NoError(t, err)
		require.NotNil(t, txInfo)

		fmt.Println("timestamp txInfo : ")
		fmt.Println(txInfo.BlockTimeStamp)

		unconfirmed, err := tronSvc.GetUnconfirmedTxAddress(context.TODO(), trxWallet.Address)
		require.NoError(t, err)
		require.Len(t, unconfirmed.Data, 1)
		require.Equal(t, hash.Hex(), "0x"+*unconfirmed.Data[0].TxID)
		require.Greater(t, unconfirmed.Data[0].NetFee, 1)
	})

}

func TestServiceTrx_GetUnconfirmedTx(t *testing.T) {
	t.Run("ShouldGetTx", func(t *testing.T) {
		// INIT
		trxSvc := trx.NewTronImpl(config.Instance())
		defer trxSvc.Close()

		// CODE UNDER TEST
		tx, err := trxSvc.GetUnconfirmedTxAddress(context.TODO(), helper.Pointer("TNjq63hm9JfqQYRRwVAtS84PRy1Ty6CU5U"))

		// EXPECTATION
		require.NoError(t, err)
		require.NotNil(t, tx)
	})
}

func TestServiceTrx_GetTx(t *testing.T) {
	t.Run("ShouldGetTx", func(t *testing.T) {
		// INIT
		trxSvc := trx.NewTronImpl(config.Instance())
		defer trxSvc.Close()

		// CODE UNDER TEST
		tx, err := trxSvc.GetTx(context.TODO(), "330c5b0f0370b6ef0a03d926cff33a14210c579dba6b195ff4735d9147649bf8")

		// EXPECTATION
		require.NoError(t, err)
		require.NotNil(t, tx)
	})
}

func TestServiceTrx_CheckAddress(t *testing.T) {
	t.Run("ShouldReturnErrorNil_WhenAddressValid", func(t *testing.T) {
		// INIT
		trxSvc := trx.NewTronImpl(config.Instance())
		defer trxSvc.Close()

		// CODE UNDER TEST
		err := trxSvc.CheckAddress("TNjq63hm9JfqQYRRwVAtS84PRy1Ty6CU5U")

		// EXPECTATION
		require.NoError(t, err)
	})

	t.Run("ShouldReturnError_WhenAddressInvalid", func(t *testing.T) {
		// INIT
		trxSvc := trx.NewTronImpl(config.Instance())
		defer trxSvc.Close()

		// CODE UNDER TEST
		err := trxSvc.CheckAddress("TNjq63hm9JfqQYRRwVKtS84PRy1Ty6CU5U")

		// EXPECTATION
		require.Error(t, err)
	})

}

func TestServiceTrx_GetCurrentBlock(t *testing.T) {
	t.Run("ShouldGetBlock", func(t *testing.T) {
		// INIT
		trxSvc := trx.NewTronImpl(config.Instance())
		defer trxSvc.Close()

		// CODE UNDER TEST
		block, err := trxSvc.GetCurrentBlock(context.TODO())

		// EXPECTATION
		require.NoError(t, err)
		require.NotNil(t, block)
		require.Greater(t, *block, int64(1))
	})
}
