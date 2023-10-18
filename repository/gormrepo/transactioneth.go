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

type ethTransaction struct {
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

func (e ethTransaction) FromModel(data model.Transaction) *ethTransaction {
	return &ethTransaction{
		Id:              data.Id,
		SenderAddress:   helper.Pointer(data.SenderAddress[0]),
		ReceiverAddress: helper.Pointer(data.ReceiverAddress[0]),
		Amount:          data.Amount,
		Fee:             data.Fee,
		Block:           data.Block,
		Confirmation:    data.Confirmation,
		Status:          data.Status,
		ReceivedAt:      data.ReceivedAt,
		UpdatedAt:       helper.Pointer(time.Now()),
	}
}

func (e ethTransaction) ToModel() *model.Transaction {
	return &model.Transaction{
		Id:              e.Id,
		SenderAddress:   []string{*e.SenderAddress},
		ReceiverAddress: []string{*e.ReceiverAddress},
		Amount:          e.Amount,
		Fee:             e.Fee,
		Block:           e.Block,
		Confirmation:    e.Confirmation,
		Status:          e.Status,
		ReceivedAt:      e.ReceivedAt,
	}
}

func (e ethTransaction) TableName() string {
	return "eth_transactions"
}

func (e *ethTransaction) BeforeCreate(db *gorm.DB) error {
	if e.Id == nil {
		return model.NewBadRequestError(helper.Pointer("id should be not nil"))
	}

	return nil
}

type EthTransactionRepo struct {
	db *gorm.DB
}

func NewEthTransactionRepository(db *gorm.DB) repository.Transaction {
	return &EthTransactionRepo{
		db: db,
	}
}

func (e *EthTransactionRepo) Add(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	gormModel := ethTransaction{}.FromModel(*transaction)

	if err := e.db.WithContext(ctx).Create(&gormModel).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, model.NewDuplicateError()
		}
		return nil, err
	}

	return gormModel.ToModel(), nil
}

func (e *EthTransactionRepo) Update(ctx context.Context, id string, transaction *model.Transaction) (*model.Transaction, error) {
	_, err := e.Get(ctx, &repository.TransactionGetFilter{Id: &id})
	if err != nil {
		return nil, err
	}

	gormModel := ethTransaction{}.FromModel(*transaction)

	tx := e.db.WithContext(ctx)
	err = tx.Model(&ethTransaction{Id: &id}).Updates(&gormModel).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, model.NewDuplicateError()
		}
		return nil, err
	}

	return e.Get(ctx, &repository.TransactionGetFilter{Id: &id})
}

func (e *EthTransactionRepo) Upsert(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	gormModel := ethTransaction{}.FromModel(*transaction)
	gormModel.UpdatedAt = helper.Pointer(time.Now())

	if err := e.db.WithContext(ctx).Table(gormModel.TableName()).Clauses(clause.OnConflict{
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

func (e *EthTransactionRepo) Get(ctx context.Context, filter *repository.TransactionGetFilter) (*model.Transaction, error) {
	transaction := ethTransaction{
		Id: filter.Id,
	}

	q := e.db.WithContext(ctx)
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
