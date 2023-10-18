package gormrepo

import (
	"time"

	"github.com/aalexanderkevin/crypto-wallet/model"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type User struct {
	Id           *string
	Email        *string
	Username     *string
	FullName     *string
	Password     *string
	PasswordSalt *string
	CreatedAt    *time.Time
}

func (u User) FromModel(data model.User) *User {
	return &User{
		Id:           data.Id,
		Email:        data.Email,
		Username:     data.Username,
		FullName:     data.FullName,
		Password:     data.Password,
		PasswordSalt: data.PasswordSalt,
		CreatedAt:    data.CreatedAt,
	}
}

func (u User) ToModel() *model.User {
	return &model.User{
		Id:           u.Id,
		Email:        u.Email,
		Username:     u.Username,
		FullName:     u.FullName,
		Password:     u.Password,
		PasswordSalt: u.PasswordSalt,
		CreatedAt:    u.CreatedAt,
	}
}

func (u User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	if u.Id == nil {
		db.Statement.SetColumn("id", ksuid.New().String())
	}

	return nil
}
