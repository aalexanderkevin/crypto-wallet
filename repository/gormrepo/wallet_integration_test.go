//go:build integration
// +build integration

package gormrepo_test

import (
	"context"
	"testing"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/helper/test"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/repository"
	"github.com/aalexanderkevin/crypto-wallet/repository/gormrepo"
	"github.com/aalexanderkevin/crypto-wallet/storage"
	"github.com/icrowley/fake"

	"github.com/stretchr/testify/require"
)

func TestWalletRepository_Get(t *testing.T) {
	t.Run("ShouldReturnNotFoundError_WhenTheIdIsNotExist", func(t *testing.T) {
		//-- init
		cfg := config.Instance()
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		//-- code under test
		walletRepo := gormrepo.NewWalletRepository(db)
		tx, err := walletRepo.Get(context.TODO(), &repository.WalletGetFilter{
			Id: helper.Pointer("invalid-id"),
		}, &cfg.Service.SeedPhraseEncryptionKey)
		require.Error(t, err)

		//-- assert
		require.EqualError(t, err, model.NewNotFoundError().Error())
		require.Nil(t, tx)
	})

	t.Run("ShouldGet_WhenTheIdExist", func(t *testing.T) {
		//-- init
		cfg := config.Instance()
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		fakeBtcTx := test.FakeWalletCreate(t, db, nil)

		//-- code under test
		walletRepo := gormrepo.NewWalletRepository(db)
		tx, err := walletRepo.Get(context.TODO(), &repository.WalletGetFilter{
			Id: fakeBtcTx.Id,
		}, &cfg.Service.SeedPhraseEncryptionKey)
		require.NoError(t, err)

		//-- assert
		require.NotNil(t, tx)
		require.Equal(t, *fakeBtcTx.Id, *tx.Id)
		require.Equal(t, fakeBtcTx.Email, tx.Email)
		require.Equal(t, fakeBtcTx.BtcAddress, tx.BtcAddress)
		require.Equal(t, *fakeBtcTx.TrxAddress, *tx.TrxAddress)
		require.Equal(t, *fakeBtcTx.EthAddress, *tx.EthAddress)
	})

}

func TestWalletRepository_Add(t *testing.T) {
	t.Run("ShouldInsertTransaction", func(t *testing.T) {
		//-- init
		cfg := config.Instance()
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		fakeWallet := test.FakeWallet(t, nil)

		//-- code under test
		walletRepo := gormrepo.NewWalletRepository(db)
		addedUser, err := walletRepo.Add(context.TODO(), &fakeWallet, &cfg.Service.SeedPhraseEncryptionKey)

		//-- assert
		require.NoError(t, err)
		require.NotNil(t, addedUser)
		require.Equal(t, fakeWallet.Id, addedUser.Id)
		require.Equal(t, fakeWallet.Email, addedUser.Email)
		require.NotEqual(t, fakeWallet.SeedPhrase, addedUser.SeedPhrase)
		require.Equal(t, fakeWallet.BtcAddress, addedUser.BtcAddress)
		require.Equal(t, fakeWallet.EthAddress, addedUser.EthAddress)
		require.Equal(t, fakeWallet.TrxAddress, addedUser.TrxAddress)
	})

	t.Run("ShouldReturnError_WhenIdAlreadyExist", func(t *testing.T) {
		//-- init
		cfg := config.Instance()
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		fakeBtcTx := test.FakeWalletCreate(t, db, nil)

		//-- code under test
		walletRepo := gormrepo.NewWalletRepository(db)
		addedUser, err := walletRepo.Add(context.TODO(), fakeBtcTx, &cfg.Service.SeedPhraseEncryptionKey)

		//-- assert
		require.Error(t, err)
		require.EqualError(t, err, model.NewDuplicateError().Error())
		require.Nil(t, addedUser)
	})

}

func TestWalletRepository_Update(t *testing.T) {
	t.Run("ShouldNotFoundError_WhenIdNotExist", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)
		invalidId := "invalid-id"

		//-- code under test
		walletRepo := gormrepo.NewWalletRepository(db)
		tx, err := walletRepo.Update(context.TODO(), invalidId, &model.Wallet{
			Email: helper.Pointer(fake.Word()),
		})
		require.Error(t, err)

		//-- assert
		require.EqualError(t, err, model.NewNotFoundError().Error())
		require.Nil(t, tx)
	})

	t.Run("ShouldUpdateUser", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		fakeTx := test.FakeWalletCreate(t, db, nil)
		updateTx := &model.Wallet{
			Email: helper.Pointer(fake.EmailAddress()),
		}

		//-- code under test
		walletRepo := gormrepo.NewWalletRepository(db)
		res, err := walletRepo.Update(context.TODO(), *fakeTx.Id, updateTx)
		require.NoError(t, err)

		//-- assert
		require.NotNil(t, res)
		require.NotEqual(t, *fakeTx.Email, *res.Email)
		require.Equal(t, *updateTx.Email, *res.Email)
		require.Equal(t, *fakeTx.BtcAddress, *res.BtcAddress)
		require.Equal(t, *fakeTx.EthAddress, *res.EthAddress)
		require.Equal(t, *fakeTx.TrxAddress, *res.TrxAddress)
		require.Equal(t, *fakeTx.SeedPhrase, *res.SeedPhrase)
	})

}
