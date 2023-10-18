package test

import (
	"context"
	"testing"
	"time"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/controller/middleware"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/repository/gormrepo"

	"github.com/golang-jwt/jwt/v5"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func FakeWallet(t *testing.T, cb func(wallet model.Wallet) model.Wallet) model.Wallet {
	t.Helper()

	fakeRp := model.Wallet{
		Id:         helper.Pointer(fake.CharactersN(7)),
		Email:      helper.Pointer(fake.EmailAddress()),
		SeedPhrase: helper.Pointer(fake.Sentence()),
		BtcAddress: helper.Pointer(fake.CharactersN(15)),
		TrxAddress: helper.Pointer(fake.CharactersN(15)),
		EthAddress: helper.Pointer(fake.CharactersN(15)),
		CreatedAt:  helper.Pointer(time.Now()),
	}
	if cb != nil {
		fakeRp = cb(fakeRp)
	}
	return fakeRp
}

func FakeWalletCreate(t *testing.T, db *gorm.DB, callback func(wallet model.Wallet) model.Wallet) *model.Wallet {
	t.Helper()
	cfg := config.Instance()

	fakeData := FakeWallet(t, callback)

	repo := gormrepo.NewWalletRepository(db)
	res, err := repo.Add(context.TODO(), &fakeData, &cfg.Service.SeedPhraseEncryptionKey)
	require.NoError(t, err)

	return res
}

func FakeTransaction(t *testing.T, cb func(transaction model.Transaction) model.Transaction) model.Transaction {
	t.Helper()

	fakeRp := model.Transaction{
		Id:              helper.Pointer(fake.CharactersN(7)),
		SenderAddress:   []string{fake.CharactersN(15)},
		ReceiverAddress: []string{fake.CharactersN(15)},
		Amount:          helper.Pointer(int64(fake.WeekdayNum())),
		Fee:             helper.Pointer(int64(fake.WeekdayNum())),
		Confirmation:    helper.Pointer(int64(fake.MonthNum())),
		Status:          helper.Pointer("pending"),
		ReceivedAt:      helper.Pointer(time.Now().Add(-2 * time.Hour)),
		CompletedAt:     helper.Pointer(time.Now()),
	}
	if cb != nil {
		fakeRp = cb(fakeRp)
	}
	return fakeRp
}

func FakeBtcTransactionCreate(t *testing.T, db *gorm.DB, callback func(transaction model.Transaction) model.Transaction) *model.Transaction {
	t.Helper()

	fakeData := FakeTransaction(t, callback)

	repo := gormrepo.NewBtcTransactionRepository(db)
	res, err := repo.Upsert(context.TODO(), &fakeData)
	require.NoError(t, err)

	return res
}

func FakeJwtToken(t *testing.T, data *string) (string, string) {
	if data == nil {
		data = helper.Pointer("email@gmail.com")
	}

	jwtClaims := middleware.JWTData{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Email: fake.EmailAddress(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	accessToken, err := token.SignedString([]byte(config.Instance().JwtSecret))
	if err != nil {
		t.Fatalf("Failed generating access token")
	}
	return accessToken, *data
}
