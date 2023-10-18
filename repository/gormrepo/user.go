package gormrepo

import (
	"context"
	"errors"

	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/repository"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.User {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) Add(ctx context.Context, user *model.User) (*model.User, error) {
	gormModel := User{}.FromModel(*user)

	if err := u.db.WithContext(ctx).Create(&gormModel).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, model.NewDuplicateError()
		}
		return nil, err
	}

	return gormModel.ToModel(), nil
}

func (u *UserRepo) Get(ctx context.Context, filter repository.UserGetFilter) (*model.User, error) {
	user := User{
		Id:    filter.Id,
		Email: filter.Email,
	}

	q := u.db.WithContext(ctx)
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

	return user.ToModel(), nil
}

func (u *UserRepo) Update(ctx context.Context, id string, user *model.User) (*model.User, error) {
	_, err := u.Get(ctx, repository.UserGetFilter{Id: &id})
	if err != nil {
		return nil, err
	}

	gormModel := User{}.FromModel(*user)

	tx := u.db.WithContext(ctx)
	err = tx.Model(&User{Id: &id}).Updates(&gormModel).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, model.NewDuplicateError()
		}
		return nil, err
	}

	return u.Get(ctx, repository.UserGetFilter{Id: &id})
}
