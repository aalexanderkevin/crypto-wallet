package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	Id           *string    `json:"id"`
	Username     *string    `json:"username"`
	Email        *string    `json:"email"`
	FullName     *string    `json:"full_name"`
	Password     *string    `json:"-"`
	PasswordSalt *string    `json:"-"`
	CreatedAt    *time.Time `json:"created_at"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(
		&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.FullName, validation.Required, validation.Length(3, 60)),
		validation.Field(&u.Password, validation.Length(6, 64)),
	)
}

func (u User) ValidateLogin() error {
	return validation.ValidateStruct(
		&u,
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.Password, validation.Required, validation.Length(6, 64)),
	)
}
