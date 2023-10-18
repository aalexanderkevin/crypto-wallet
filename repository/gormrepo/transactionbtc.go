package gormrepo

import (
	"context"
	"errors"
	"time"

	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/repository"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type btcTransaction struct {
	Id              *string
	SenderAddress   pq.StringArray `gorm:"type:text[]"`
	ReceiverAddress pq.StringArray `gorm:"type:text[]"`
	Amount          *int64
	Fee             *int64
	Confirmation    *int64
	Status          *string
	ReceivedAt      *time.Time
	CompletedAt     *time.Time
}

func (b btcTransaction) FromModel(data model.Transaction) *btcTransaction {
	return &btcTransaction{
		Id:              data.Id,
		SenderAddress:   data.SenderAddress,
		ReceiverAddress: data.ReceiverAddress,
		Amount:          data.Amount,
		Fee:             data.Fee,
		Confirmation:    data.Confirmation,
		Status:          data.Status,
		ReceivedAt:      data.ReceivedAt,
		CompletedAt:     data.CompletedAt,
	}
}

func (b btcTransaction) ToModel() *model.Transaction {
	return &model.Transaction{
		Id:              b.Id,
		SenderAddress:   b.SenderAddress,
		ReceiverAddress: b.ReceiverAddress,
		Amount:          b.Amount,
		Fee:             b.Fee,
		Confirmation:    b.Confirmation,
		Status:          b.Status,
		ReceivedAt:      b.ReceivedAt,
		CompletedAt:     b.CompletedAt,
	}
}

func (b btcTransaction) TableName() string {
	return "btc_transactions"
}

func (b *btcTransaction) BeforeCreate(db *gorm.DB) error {
	if b.Id == nil {
		return model.NewBadRequestError(helper.Pointer("id should be not nil"))
	}

	return nil
}

type BtcTransactionRepo struct {
	db *gorm.DB
}

func NewBtcTransactionRepository(db *gorm.DB) repository.Transaction {
	return &BtcTransactionRepo{
		db: db,
	}
}

func (b *BtcTransactionRepo) Add(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	gormModel := btcTransaction{}.FromModel(*transaction)

	if err := b.db.WithContext(ctx).Create(&gormModel).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, model.NewDuplicateError()
		}
		return nil, err
	}

	return gormModel.ToModel(), nil
}

func (b *BtcTransactionRepo) Update(ctx context.Context, id string, transaction *model.Transaction) (*model.Transaction, error) {
	_, err := b.Get(ctx, &repository.TransactionGetFilter{Id: &id})
	if err != nil {
		return nil, err
	}

	gormModel := btcTransaction{}.FromModel(*transaction)

	tx := b.db.WithContext(ctx)
	err = tx.Model(&btcTransaction{Id: &id}).Updates(&gormModel).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, model.NewDuplicateError()
		}
		return nil, err
	}

	return b.Get(ctx, &repository.TransactionGetFilter{Id: &id})
}

func (b *BtcTransactionRepo) Upsert(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	gormModel := btcTransaction{}.FromModel(*transaction)

	if err := b.db.WithContext(ctx).Table(gormModel.TableName()).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"sender_address", "receiver_address", "amount", "fee", "confirmation", "status", "received_at", "completed_at"}),
	}).Create(&gormModel).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, model.NewDuplicateError()
		}
		return nil, err
	}

	return gormModel.ToModel(), nil
}

func (b *BtcTransactionRepo) Get(ctx context.Context, filter *repository.TransactionGetFilter) (*model.Transaction, error) {
	transaction := btcTransaction{
		Id: filter.Id,
	}

	q := b.db.WithContext(ctx)
	if filter.Id != nil {
		q = q.Where("id = ?", filter.Id)
	}

	if filter.SenderAddress != nil {
		q = q.Where("sender_address = ?", filter.SenderAddress)
	}

	if filter.ReceiverAddress != nil {
		q = q.Where("receiver_address = ?", filter.ReceiverAddress)
	}

	if filter.Status != nil {
		q = q.Where("status = ?", filter.Status)
	}

	err := q.First(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.NewNotFoundError()
		}
		return nil, err
	}

	return transaction.ToModel(), nil
}
