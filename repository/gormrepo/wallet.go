package gormrepo

import (
	"context"
	"errors"
	"time"

	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/repository"
	"github.com/segmentio/ksuid"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type WalletRepo struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) repository.Wallet {
	return &WalletRepo{
		db: db,
	}
}

type Wallet struct {
	Id         *string
	Email      *string
	SeedPhrase []byte
	BtcAddress *string
	EthAddress *string
	TrxAddress *string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

func (w Wallet) FromModel(data *model.Wallet, key *string) (wallet *Wallet, err error) {
	var seedPhrase []byte
	if data.SeedPhrase != nil {
		seedPhrase = []byte(*data.SeedPhrase)
		if key != nil {
			seedPhrase, err = helper.EncryptSeedPhrase(*data.SeedPhrase, *key)
			if err != nil {
				return nil, err
			}
		}
	}

	return &Wallet{
		Id:         data.Id,
		Email:      data.Email,
		SeedPhrase: seedPhrase,
		BtcAddress: data.BtcAddress,
		EthAddress: data.EthAddress,
		TrxAddress: data.TrxAddress,
		CreatedAt:  data.CreatedAt,
		UpdatedAt:  data.UpdatedAt,
	}, nil
}

func (w Wallet) ToModel(key *string) (wallet *model.Wallet, err error) {
	seedPhrase := w.SeedPhrase
	if key != nil {
		seedPhrase, err = helper.DecryptSeedPhrase(w.SeedPhrase, *key)
		if err != nil {
			return nil, err
		}
	}

	return &model.Wallet{
		Id:         w.Id,
		Email:      w.Email,
		SeedPhrase: helper.Pointer(string(seedPhrase)),
		BtcAddress: w.BtcAddress,
		EthAddress: w.EthAddress,
		TrxAddress: w.TrxAddress,
		CreatedAt:  w.CreatedAt,
		UpdatedAt:  w.UpdatedAt,
	}, nil
}

func (w Wallet) TableName() string {
	return "wallets"
}

func (w *Wallet) BeforeCreate(db *gorm.DB) error {
	if w.Id == nil {
		db.Statement.SetColumn("id", ksuid.New().String())
	}

	return nil
}

func (w *WalletRepo) Add(ctx context.Context, wallet *model.Wallet, encryptionKey *string) (*model.Wallet, error) {
	gormModel, err := Wallet{}.FromModel(wallet, encryptionKey)
	if err != nil {
		return nil, err
	}

	if err := w.db.WithContext(ctx).Create(&gormModel).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, model.NewDuplicateError()
		}
		return nil, err
	}

	return gormModel.ToModel(nil)
}

func (w *WalletRepo) Get(ctx context.Context, filter *repository.WalletGetFilter, encryptionKey *string) (*model.Wallet, error) {
	user := Wallet{
		Id:    filter.Id,
		Email: filter.Email,
	}

	q := w.db.WithContext(ctx)
	if filter.Id != nil {
		q = q.Where("id = ?", filter.Id)
	}
	if filter.Email != nil {
		q = q.Where("email = ?", filter.Email)
	}

	err := q.First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.NewNotFoundError()
		}
		return nil, err
	}

	return user.ToModel(encryptionKey)
}

func (w *WalletRepo) Update(ctx context.Context, id string, wallet *model.Wallet) (*model.Wallet, error) {
	_, err := w.Get(ctx, &repository.WalletGetFilter{Id: &id}, nil)
	if err != nil {
		return nil, err
	}

	gormModel, err := Wallet{}.FromModel(wallet, nil)
	if err != nil {
		return nil, err
	}

	tx := w.db.WithContext(ctx)
	err = tx.Model(&Wallet{Id: &id}).Updates(&gormModel).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, model.NewDuplicateError()
		}
		return nil, err
	}

	return w.Get(ctx, &repository.WalletGetFilter{Id: &id}, nil)
}
