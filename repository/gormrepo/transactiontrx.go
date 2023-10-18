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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type trxTransaction struct {
	Id              *string
	SenderAddress   *string
	ReceiverAddress *string
	Amount          *int64
	Fee             *int64
	Block           *int64
	Confirmation    *int64
	Status          *string
	ReceivedAt      *time.Time
	UpdatedAt       *time.Time
}

func (t trxTransaction) FromModel(data model.Transaction) *trxTransaction {
	return &trxTransaction{
		Id:              data.Id,
		SenderAddress:   helper.Pointer(data.SenderAddress[0]),
		ReceiverAddress: helper.Pointer(data.ReceiverAddress[0]),
		Amount:          data.Amount,
		Fee:             data.Fee,
		Block:           data.Block,
		Confirmation:    data.Confirmation,
		Status:          data.Status,
		ReceivedAt:      data.ReceivedAt,
	}
}

func (t trxTransaction) ToModel() *model.Transaction {
	return &model.Transaction{
		Id:              t.Id,
		SenderAddress:   []string{*t.SenderAddress},
		ReceiverAddress: []string{*t.ReceiverAddress},
		Amount:          t.Amount,
		Fee:             t.Fee,
		Block:           t.Block,
		Confirmation:    t.Confirmation,
		Status:          t.Status,
		ReceivedAt:      t.ReceivedAt,
	}
}

func (t trxTransaction) TableName() string {
	return "trx_transactions"
}

func (t *trxTransaction) BeforeCreate(db *gorm.DB) error {
	if t.Id == nil {
		return model.NewBadRequestError(helper.Pointer("id should be not nil"))
	}

	return nil
}

type TrxTransactionRepo struct {
	db *gorm.DB
}

func NewTrxTransactionRepository(db *gorm.DB) repository.Transaction {
	return &TrxTransactionRepo{
		db: db,
	}
}

func (t *TrxTransactionRepo) Add(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	gormModel := trxTransaction{}.FromModel(*transaction)

	if err := t.db.WithContext(ctx).Create(&gormModel).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, model.NewDuplicateError()
		}
		return nil, err
	}

	return gormModel.ToModel(), nil
}

func (t *TrxTransactionRepo) Update(ctx context.Context, id string, transaction *model.Transaction) (*model.Transaction, error) {
	_, err := t.Get(ctx, &repository.TransactionGetFilter{Id: &id})
	if err != nil {
		return nil, err
	}

	gormModel := trxTransaction{}.FromModel(*transaction)

	tx := t.db.WithContext(ctx)
	err = tx.Model(&trxTransaction{Id: &id}).Updates(&gormModel).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, model.NewDuplicateError()
		}
		return nil, err
	}

	return t.Get(ctx, &repository.TransactionGetFilter{Id: &id})
}

func (t *TrxTransactionRepo) Upsert(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	gormModel := trxTransaction{}.FromModel(*transaction)
	gormModel.UpdatedAt = helper.Pointer(time.Now())

	if err := t.db.WithContext(ctx).Table(gormModel.TableName()).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"sender_address", "receiver_address", "amount", "fee", "block", "confirmation", "status", "received_at", "updated_at"}),
	}).Create(&gormModel).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, model.NewDuplicateError()
		}
		return nil, err
	}

	return gormModel.ToModel(), nil
}

func (t *TrxTransactionRepo) Get(ctx context.Context, filter *repository.TransactionGetFilter) (*model.Transaction, error) {
	transaction := trxTransaction{
		Id: filter.Id,
	}

	q := t.db.WithContext(ctx)
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
