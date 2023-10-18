package usecase

import (
	"context"

	"github.com/aalexanderkevin/crypto-wallet/container"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/repository"

	"github.com/segmentio/ksuid"
)

type User struct {
	repository.User
}

func NewUser(u *container.Container) *User {
	return &User{
		User: u.UserRepo(),
	}
}

func (u *User) Register(ctx context.Context, req *model.User) (*model.User, error) {
	logger := helper.GetLogger(ctx).WithField("method", "usecase.User.Register")

	if err := req.Validate(); err != nil {
		logger.WithError(err).Warning("Not Valid Request")
		return nil, model.NewParameterError(helper.Pointer(err.Error()))
	}

	_, err := u.User.Get(ctx, repository.UserGetFilter{
		Email: req.Email,
	})
	if err == nil {
		err = model.NewDuplicateError()
		logger.WithError(err).Warning("Email already exist")
		return nil, err
	} else if !model.IsNotFoundError(err) {
		logger.WithError(err).Warning("Failed to check existing email")
		return nil, err
	}

	passwordSalt := ksuid.New().String()
	req.Password = helper.Pointer(helper.Hash(passwordSalt, *req.Password))
	req.PasswordSalt = helper.Pointer(passwordSalt)

	res, err := u.User.Add(ctx, req)
	if err != nil {
		logger.WithError(err).Warning("Failed insert User")
		return nil, err
	}
	res.Password = nil
	res.PasswordSalt = nil

	return res, nil
}

func (u *User) Login(ctx context.Context, req *model.User) (*model.User, error) {
	logger := helper.GetLogger(ctx).WithField("method", "usecase.User.Login")

	if err := req.ValidateLogin(); err != nil {
		logger.WithError(err).Warning("Not Valid Request")
		return nil, model.NewParameterError(helper.Pointer(err.Error()))
	}

	user, err := u.User.Get(ctx, repository.UserGetFilter{
		Username: req.Username,
	})
	if err != nil {
		logger.WithError(err).Warning("Failed get User")
		return nil, err
	}

	reqPassword := helper.Hash(*user.PasswordSalt, *req.Password)
	if reqPassword != *user.Password {
		err := model.NewInvalidPasswordError()
		logger.WithError(err).Warning("Invalid password")
		return nil, err
	}

	return user, nil
}

func (u *User) Update(ctx context.Context, email *string, req *model.User) (*model.User, error) {
	logger := helper.GetLogger(ctx).WithField("method", "usecase.Update")

	if email == nil {
		logger.Error("missing email")
		return nil, model.NewParameterError(helper.Pointer("missing email"))
	}

	user, err := u.User.Get(ctx, repository.UserGetFilter{
		Email: email,
	})
	if err != nil {
		logger.WithError(err).Warning("Failed get User")
		return nil, err
	}

	res, err := u.User.Update(ctx, *user.Id, req)
	if err != nil {
		logger.WithError(err).Warning("Failed update User")
		return nil, err
	}

	return res, nil
}
