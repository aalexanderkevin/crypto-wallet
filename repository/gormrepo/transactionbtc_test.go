//go:build integration
// +build integration

package gormrepo_test

import (
	"context"
	"testing"

	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/helper/test"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/repository"
	"github.com/aalexanderkevin/crypto-wallet/repository/gormrepo"
	"github.com/aalexanderkevin/crypto-wallet/storage"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
)

func TestBtcTransactionRepository_Upsert(t *testing.T) {
	t.Run("ShouldInsertBtcTransaction", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		fakeTransaction := test.FakeTransaction(t, nil)

		//-- code under test
		btcTxRepo := gormrepo.NewBtcTransactionRepository(db)
		addedUser, err := btcTxRepo.Upsert(context.TODO(), &fakeTransaction)

		//-- assert
		require.NoError(t, err)
		require.NotNil(t, addedUser)
		require.Equal(t, fakeTransaction.Id, addedUser.Id)
		require.Equal(t, fakeTransaction.SenderAddress, addedUser.SenderAddress)
		require.Equal(t, fakeTransaction.ReceivedAt, addedUser.ReceivedAt)
		require.Equal(t, fakeTransaction.Fee, addedUser.Fee)
		require.Equal(t, fakeTransaction.Amount, addedUser.Amount)
		require.Equal(t, fakeTransaction.Confirmation, addedUser.Confirmation)
		require.Equal(t, fakeTransaction.Status, addedUser.Status)
		require.Equal(t, fakeTransaction.ReceivedAt, addedUser.ReceivedAt)
		require.Equal(t, fakeTransaction.CompletedAt, addedUser.CompletedAt)
	})

	t.Run("ShouldUpdate_WhenTheIdAlreadyExist", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		fakeBtcTx := test.FakeBtcTransactionCreate(t, db, func(transaction model.Transaction) model.Transaction {
			transaction.Confirmation = helper.Pointer[int64](2)
			return transaction
		})

		fakeBtcTx.Confirmation = helper.Pointer[int64](4)

		//-- code under test
		btcTxRepo := gormrepo.NewBtcTransactionRepository(db)
		updatedUser, err := btcTxRepo.Upsert(context.TODO(), fakeBtcTx)
		require.NoError(t, err)

		res, err := btcTxRepo.Get(context.TODO(), &repository.TransactionGetFilter{
			Id: fakeBtcTx.Id,
		})

		//-- assert
		require.NoError(t, err)
		require.Equal(t, fakeBtcTx.Id, updatedUser.Id)
		require.Equal(t, fakeBtcTx.SenderAddress, res.SenderAddress)
		require.Equal(t, fakeBtcTx.ReceiverAddress, res.ReceiverAddress)
		require.Equal(t, fakeBtcTx.Confirmation, res.Confirmation)
	})

}

func TestBtcTransactionRepository_Get(t *testing.T) {
	t.Run("ShouldReturnNotFoundError_WhenTheIdIsNotExist", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		//-- code under test
		btcTxRepo := gormrepo.NewBtcTransactionRepository(db)
		tx, err := btcTxRepo.Get(context.TODO(), &repository.TransactionGetFilter{
			Id: helper.Pointer("invalid-id"),
		})
		require.Error(t, err)

		//-- assert
		require.EqualError(t, err, model.NewNotFoundError().Error())
		require.Nil(t, tx)
	})

	t.Run("ShouldGet_WhenTheIdExist", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		fakeBtcTx := test.FakeBtcTransactionCreate(t, db, nil)

		//-- code under test
		btcTxRepo := gormrepo.NewBtcTransactionRepository(db)
		tx, err := btcTxRepo.Get(context.TODO(), &repository.TransactionGetFilter{
			Id: fakeBtcTx.Id,
		})
		require.NoError(t, err)

		//-- assert
		require.NotNil(t, tx)
		require.Equal(t, *fakeBtcTx.Id, *tx.Id)
		require.Equal(t, fakeBtcTx.ReceiverAddress, tx.ReceiverAddress)
		require.Equal(t, fakeBtcTx.SenderAddress, tx.SenderAddress)
		require.Equal(t, *fakeBtcTx.Amount, *tx.Amount)
		require.Equal(t, *fakeBtcTx.Fee, *tx.Fee)
		require.Equal(t, *fakeBtcTx.Confirmation, *tx.Confirmation)
		require.Equal(t, *fakeBtcTx.Status, *tx.Status)
	})

}

func TestBtcTransactionRepository_Add(t *testing.T) {
	t.Run("ShouldInsertTransaction", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		fakeTransaction := test.FakeTransaction(t, nil)

		//-- code under test
		btcTxRepo := gormrepo.NewBtcTransactionRepository(db)
		addedUser, err := btcTxRepo.Add(context.TODO(), &fakeTransaction)

		//-- assert
		require.NoError(t, err)
		require.NotNil(t, addedUser)
		require.Equal(t, fakeTransaction.Id, addedUser.Id)
		require.Equal(t, fakeTransaction.SenderAddress, addedUser.SenderAddress)
		require.Equal(t, fakeTransaction.ReceiverAddress, addedUser.ReceiverAddress)
		require.Equal(t, fakeTransaction.Amount, addedUser.Amount)
		require.Equal(t, fakeTransaction.Fee, addedUser.Fee)
	})

	t.Run("ShouldReturnError_WhenIdAlreadyExist", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		fakeBtcTx := test.FakeBtcTransactionCreate(t, db, nil)

		//-- code under test
		btcTxRepo := gormrepo.NewBtcTransactionRepository(db)
		addedUser, err := btcTxRepo.Add(context.TODO(), fakeBtcTx)

		//-- assert
		require.Error(t, err)
		require.EqualError(t, err, model.NewDuplicateError().Error())
		require.Nil(t, addedUser)
	})

	t.Run("ShouldReturnError_WhenIdNotProvided", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		fakeTransaction := test.FakeTransaction(t, nil)
		fakeTransaction.Id = nil

		//-- code under test
		btcTxRepo := gormrepo.NewBtcTransactionRepository(db)
		addedUser, err := btcTxRepo.Add(context.TODO(), &fakeTransaction)

		//-- assert
		require.Error(t, err)
		require.EqualError(t, err, model.NewBadRequestError(helper.Pointer("id should be not nil")).Error())
		require.Nil(t, addedUser)
	})

}

func TestBtcTransactionRepository_Update(t *testing.T) {
	t.Run("ShouldNotFoundError_WhenIdNotExist", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)
		invalidId := "invalid-id"

		//-- code under test
		btcTxRepo := gormrepo.NewBtcTransactionRepository(db)
		tx, err := btcTxRepo.Update(context.TODO(), invalidId, &model.Transaction{
			ReceiverAddress: []string{fake.Word()},
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

		fakeTx := test.FakeBtcTransactionCreate(t, db, nil)
		updateTx := &model.Transaction{
			Status: helper.Pointer("success"),
		}

		//-- code under test
		btcTxRepo := gormrepo.NewBtcTransactionRepository(db)
		res, err := btcTxRepo.Update(context.TODO(), *fakeTx.Id, updateTx)
		require.NoError(t, err)

		//-- assert
		require.NotNil(t, res)
		require.NotEqual(t, *fakeTx.Status, *res.Status)
		require.Equal(t, *updateTx.Status, *res.Status)
		require.Equal(t, fakeTx.ReceiverAddress, res.ReceiverAddress)
		require.Equal(t, fakeTx.SenderAddress, res.SenderAddress)
		require.Equal(t, *fakeTx.Amount, *res.Amount)
		require.Equal(t, *fakeTx.Fee, *res.Fee)
	})

}
