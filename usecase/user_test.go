package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aalexanderkevin/crypto-wallet/container"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/helper/test"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/repository"
	"github.com/aalexanderkevin/crypto-wallet/repository/mocks"
	"github.com/aalexanderkevin/crypto-wallet/usecase"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUser_Register(t *testing.T) {
	t.Parallel()
	t.Run("ShouldReturnError_WhenEmailIsMissing", func(t *testing.T) {
		t.Parallel()
		// INIT
		appContainer := container.Container{}

		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = nil
			return user
		})

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, err := uc.Register(context.Background(), &fakeUser)
		require.Error(t, err)
		require.True(t, model.IsParameterError(err))
		require.Nil(t, res)
	})

	t.Run("ShouldReturnErrorDuplicate_WhenEmailAlreadyExist", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("fakeemail@gmail.com")
			return user
		})

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Email: fakeUser.Email,
		}).Return(&fakeUser, nil).Once()

		appContainer := container.Container{}
		appContainer.SetUserRepo(userMock)

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, err := uc.Register(context.Background(), &fakeUser)
		require.Error(t, err)
		require.EqualError(t, err, model.NewDuplicateError().Error())
		require.Nil(t, res)

		userMock.AssertExpectations(t)
	})

	t.Run("ShouldReturnError_WhenErrorGetUser", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("fakeemail@gmail.com")
			return user
		})
		errMock := errors.New("error get")

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Email: fakeUser.Email,
		}).Return(nil, errMock).Once()

		appContainer := container.Container{}
		appContainer.SetUserRepo(userMock)

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, err := uc.Register(context.Background(), &fakeUser)
		require.Error(t, err)
		require.EqualError(t, err, errMock.Error())
		require.Nil(t, res)

		userMock.AssertExpectations(t)
	})

	t.Run("ShouldReturnError_WhenErrorAddUser", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("fakeemail@gmail.com")
			user.PasswordSalt = nil
			return user
		})
		password := *fakeUser.Password

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Email: fakeUser.Email,
		}).Return(nil, model.NewNotFoundError()).Once()
		userMock.On("Add", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			require.Equal(t, *fakeUser.Email, *u.Email)
			require.Equal(t, *fakeUser.FullName, *u.FullName)

			password := helper.Pointer(helper.Hash(*u.PasswordSalt, password))
			require.Equal(t, *password, *u.Password)
			return true
		})).Return(nil, errors.New("error insert")).Once()

		appContainer := container.Container{}
		appContainer.SetUserRepo(userMock)

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, err := uc.Register(context.Background(), &fakeUser)
		require.Error(t, err)
		require.EqualError(t, err, "error insert")
		require.Nil(t, res)

		userMock.AssertExpectations(t)
	})

	t.Run("ShouldRegisterNewUser", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("fakeemail@gmail.com")
			user.PasswordSalt = nil
			return user
		})
		cleanPassword := *fakeUser.Password

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Email: fakeUser.Email,
		}).Return(nil, model.NewNotFoundError()).Once()
		userMock.On("Add", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			require.Equal(t, *fakeUser.Email, *u.Email)
			require.Equal(t, *fakeUser.FullName, *u.FullName)

			password := helper.Pointer(helper.Hash(*u.PasswordSalt, cleanPassword))
			require.Equal(t, *password, *u.Password)
			return true
		})).Return(&model.User{
			Id:           helper.Pointer(fake.CharactersN(6)),
			Email:        fakeUser.Email,
			FullName:     fakeUser.FullName,
			Password:     fakeUser.Password,
			PasswordSalt: fakeUser.PasswordSalt,
			CreatedAt:    fakeUser.CreatedAt,
		}, nil).Once()

		appContainer := container.Container{}
		appContainer.SetUserRepo(userMock)

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, err := uc.Register(context.Background(), &fakeUser)
		require.NoError(t, err)
		require.NotNil(t, res.Id)
		require.Equal(t, *fakeUser.Email, *res.Email)
		require.Equal(t, *fakeUser.FullName, *res.FullName)
		require.Nil(t, res.Password)
		require.Nil(t, res.PasswordSalt)

		userMock.AssertExpectations(t)
	})
}

func TestUser_Login(t *testing.T) {
	t.Parallel()
	t.Run("ShouldReturnError_WhenUserIsMissing", func(t *testing.T) {
		t.Parallel()
		// INIT
		appContainer := container.Container{}

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, err := uc.Login(context.Background(), &model.User{
			Password: helper.Pointer("password"),
		})
		require.Error(t, err)
		require.True(t, model.IsParameterError(err))
		require.Nil(t, res)

	})

	t.Run("ShouldReturnError_WhenErrorGetUser", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, nil)

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Username: fakeUser.Username,
		}).Return(nil, errors.New("error get")).Once()

		appContainer := container.Container{}
		appContainer.SetUserRepo(userMock)

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, err := uc.Login(context.Background(), &model.User{
			Username: fakeUser.Username,
			Password: fakeUser.Password,
		})
		require.Error(t, err)
		require.EqualError(t, err, "error get")
		require.Nil(t, res)

		userMock.AssertExpectations(t)
	})

	t.Run("ShouldReturnError_WhenPasswordIsIncorrect", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, nil)

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Username: fakeUser.Username,
		}).Return(&fakeUser, nil).Once()

		appContainer := container.Container{}
		appContainer.SetUserRepo(userMock)

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, err := uc.Login(context.Background(), &model.User{
			Username: fakeUser.Username,
			Password: fakeUser.Password,
		})
		require.Error(t, err)
		require.EqualError(t, err, model.NewBadRequestError(helper.Pointer("invalid password")).Error())
		require.Nil(t, res)

		userMock.AssertExpectations(t)
	})

	t.Run("ShouldLoginSuccess_WhenEmailAndPasswordAreCorrect", func(t *testing.T) {
		t.Parallel()
		// INIT
		password := fake.CharactersN(7)
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.PasswordSalt = helper.Pointer(fake.CharactersN(7))
			user.Password = helper.Pointer(helper.Hash(*user.PasswordSalt, password))
			return user
		})

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Username: fakeUser.Username,
		}).Return(&fakeUser, nil).Once()

		appContainer := container.Container{}
		appContainer.SetUserRepo(userMock)

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, err := uc.Login(context.Background(), &model.User{
			Username: fakeUser.Username,
			Password: &password,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, *fakeUser.Email, *res.Email)
		require.Equal(t, *fakeUser.FullName, *res.FullName)
		require.Equal(t, *fakeUser.Password, *res.Password)
		require.Equal(t, *fakeUser.PasswordSalt, *res.PasswordSalt)

		userMock.AssertExpectations(t)
	})
}

func TestUser_Update(t *testing.T) {
	t.Parallel()
	t.Run("ShouldReturnError_WhenEmailIsMissing", func(t *testing.T) {
		t.Parallel()
		// INIT
		appContainer := container.Container{}

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, err := uc.Update(context.Background(), nil, &model.User{})
		require.Error(t, err)
		require.True(t, model.IsParameterError(err))
		require.Nil(t, res)

	})

	t.Run("ShouldReturnError_WhenErrorGetUser", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("invalidEmail@gmail.com")
			return user
		})
		updateUser := &model.User{
			FullName: helper.Pointer("new-fullname"),
		}
		errMock := errors.New("error get")

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Email: fakeUser.Email,
		}).Return(nil, errMock).Once()

		appContainer := container.Container{}
		appContainer.SetUserRepo(userMock)

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, err := uc.Update(context.Background(), fakeUser.Email, updateUser)
		require.Error(t, err)
		require.EqualError(t, err, errMock.Error())
		require.Nil(t, res)

		userMock.AssertExpectations(t)
	})

	t.Run("ShouldReturnError_WhenErrorUpdateUser", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("invalidEmail@gmail.com")
			return user
		})
		updateUser := &model.User{
			FullName: helper.Pointer("new-fullname"),
		}

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Email: fakeUser.Email,
		}).Return(&fakeUser, nil).Once()
		userMock.On("Update", mock.Anything, *fakeUser.Id, updateUser).Return(nil, errors.New("error update")).Once()

		appContainer := container.Container{}
		appContainer.SetUserRepo(userMock)

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, err := uc.Update(context.Background(), fakeUser.Email, updateUser)
		require.Error(t, err)
		require.EqualError(t, err, "error update")
		require.Nil(t, res)

		userMock.AssertExpectations(t)
	})

	t.Run("ShouldUpdateUser", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("valid@gmail.com")
			return user
		})
		updateUser := &model.User{
			Email:    helper.Pointer(fake.EmailAddress()),
			FullName: helper.Pointer(fake.FullName()),
		}
		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Email: fakeUser.Email,
		}).Return(&fakeUser, nil).Once()
		userMock.On("Update", mock.Anything, *fakeUser.Id, updateUser).Return(updateUser, nil).Once()

		appContainer := container.Container{}
		appContainer.SetUserRepo(userMock)

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, err := uc.Update(context.Background(), fakeUser.Email, updateUser)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, updateUser.Email, res.Email)
		require.Equal(t, updateUser.FullName, res.FullName)

		userMock.AssertExpectations(t)
	})
}
