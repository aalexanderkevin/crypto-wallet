package handler_test

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
	cegrpc "github.com/aalexanderkevin/crypto-wallet/transport/grpc/crypto-wallet"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUser_Register(t *testing.T) {
	t.Parallel()
	t.Run("ShouldReturnErrorInvalidArgument_WhenEmailIsMissing", func(t *testing.T) {
		t.Parallel()
		// INIT
		ctx := context.Background()
		fakeUser := test.FakeUser(t, nil)
		req := &cegrpc.RegisterRequest{
			Username: *fakeUser.Username,
			Fullname: *fakeUser.FullName,
			Password: *fakeUser.Password,
		}

		conn, closeResourcesFn := SetupGRPCConn(t, ctx, nil)
		client := cegrpc.NewCryptoWalletClient(conn)
		defer closeResourcesFn()

		// CODE UNDER TEST
		response, err := performGRPCRequest(ctx, client, "Register", req)

		// EXPECTATION
		require.Error(t, err)
		require.Nil(t, response)

		status, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.InvalidArgument, status.Code())
	})

	t.Run("ShouldReturnError_WhenFailedAddNewUser", func(t *testing.T) {
		t.Parallel()
		// INIT
		ctx := context.Background()
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("email@gmail.com")
			return user
		})
		req := &cegrpc.RegisterRequest{
			Username: *fakeUser.Username,
			Email:    *fakeUser.Email,
			Fullname: *fakeUser.FullName,
			Password: *fakeUser.Password,
		}

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Email: fakeUser.Email,
		}).Return(nil, model.NewNotFoundError()).Once()
		userMock.On("Add", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			require.Equal(t, *fakeUser.Email, *u.Email)
			require.Equal(t, *fakeUser.FullName, *u.FullName)

			password := helper.Pointer(helper.Hash(*u.PasswordSalt, *fakeUser.Password))
			require.Equal(t, *password, *u.Password)
			return true
		})).Return(nil, errors.New("error insert")).Once()

		conn, closeResourcesFn := SetupGRPCConn(t, ctx, func(appContainer *container.Container) *container.Container {
			appContainer.SetUserRepo(userMock)
			return appContainer
		})
		client := cegrpc.NewCryptoWalletClient(conn)
		defer closeResourcesFn()

		// CODE UNDER TEST
		response, err := performGRPCRequest(ctx, client, "Register", req)

		// EXPECTATION
		require.Error(t, err)
		require.Nil(t, response)

		status, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Internal, status.Code())
		require.Equal(t, "error insert", status.Message())

		userMock.AssertExpectations(t)
	})

	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		t.Parallel()
		// INIT
		ctx := context.Background()
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("email@gmail.com")
			return user
		})
		req := &cegrpc.RegisterRequest{
			Username: *fakeUser.Username,
			Fullname: *fakeUser.FullName,
			Email:    *fakeUser.Email,
			Password: *fakeUser.Password,
		}

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Email: fakeUser.Email,
		}).Return(nil, model.NewNotFoundError()).Once()
		userMock.On("Add", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			require.Equal(t, *fakeUser.Email, *u.Email)
			require.Equal(t, *fakeUser.FullName, *u.FullName)

			password := helper.Pointer(helper.Hash(*u.PasswordSalt, *fakeUser.Password))
			require.Equal(t, *password, *u.Password)
			return true
		})).Return(&fakeUser, nil).Once()

		conn, closeResourcesFn := SetupGRPCConn(t, ctx, func(appContainer *container.Container) *container.Container {
			appContainer.SetUserRepo(userMock)
			return appContainer
		})
		client := cegrpc.NewCryptoWalletClient(conn)
		defer closeResourcesFn()

		// CODE UNDER TEST
		response, err := performGRPCRequest(ctx, client, "Register", req)

		// EXPECTATION
		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, "Register successful", response.RegisterResponse.Message)
		require.True(t, response.RegisterResponse.Success)

		userMock.AssertExpectations(t)
	})

}

func TestUser_Login(t *testing.T) {
	t.Run("ShouldReturnInvalidArgument_WhenUsernameIsMissing", func(t *testing.T) {
		// INIT
		ctx := context.Background()
		fakeUser := test.FakeUser(t, nil)
		req := &cegrpc.LoginRequest{
			Username: "",
			Password: *fakeUser.Password,
		}

		conn, closeResourcesFn := SetupGRPCConn(t, ctx, nil)
		client := cegrpc.NewCryptoWalletClient(conn)
		defer closeResourcesFn()

		// CODE UNDER TEST
		response, err := performGRPCRequest(ctx, client, "Login", req)

		// EXPECTATION
		require.Error(t, err)
		require.Nil(t, response)

		status, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.InvalidArgument, status.Code())
	})

	t.Run("ShouldReturnError_WhenErrorGetUser", func(t *testing.T) {
		// INIT
		ctx := context.Background()
		fakeUser := test.FakeUser(t, nil)
		req := &cegrpc.LoginRequest{
			Username: *fakeUser.Username,
			Password: *fakeUser.Password,
		}

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Username: fakeUser.Username,
		}).Return(nil, errors.New("error get")).Once()

		conn, closeResourcesFn := SetupGRPCConn(t, ctx, func(appContainer *container.Container) *container.Container {
			appContainer.SetUserRepo(userMock)
			return appContainer
		})
		client := cegrpc.NewCryptoWalletClient(conn)
		defer closeResourcesFn()

		// CODE UNDER TEST
		response, err := performGRPCRequest(ctx, client, "Login", req)

		// EXPECTATION
		require.Error(t, err)
		require.Nil(t, response)

		status, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Internal, status.Code())
		require.Equal(t, "error get", status.Message())

		userMock.AssertExpectations(t)
	})

	t.Run("ShouldReturnError_WhenPasswordIsIncorrect", func(t *testing.T) {
		// INIT
		ctx := context.Background()
		fakeUser := test.FakeUser(t, nil)
		req := &cegrpc.LoginRequest{
			Username: *fakeUser.Username,
			Password: *fakeUser.Password,
		}

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Username: fakeUser.Username,
		}).Return(&fakeUser, nil).Once()

		conn, closeResourcesFn := SetupGRPCConn(t, ctx, func(appContainer *container.Container) *container.Container {
			appContainer.SetUserRepo(userMock)
			return appContainer
		})
		client := cegrpc.NewCryptoWalletClient(conn)
		defer closeResourcesFn()

		// CODE UNDER TEST
		response, err := performGRPCRequest(ctx, client, "Login", req)

		// EXPECTATION
		require.Error(t, err)
		require.Nil(t, response)

		status, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Unauthenticated, status.Code())

		userMock.AssertExpectations(t)
	})

	t.Run("ShouldLoginSuccess_WhenEmailAndPasswordAreCorrect", func(t *testing.T) {
		// INIT
		ctx := context.Background()
		password := fake.CharactersN(7)
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.PasswordSalt = helper.Pointer(fake.CharactersN(7))
			user.Password = helper.Pointer(helper.Hash(*user.PasswordSalt, password))
			return user
		})
		req := &cegrpc.LoginRequest{
			Username: *fakeUser.Username,
			Password: password,
		}

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Username: fakeUser.Username,
		}).Return(&fakeUser, nil).Once()

		conn, closeResourcesFn := SetupGRPCConn(t, ctx, func(appContainer *container.Container) *container.Container {
			appContainer.SetUserRepo(userMock)
			return appContainer
		})
		client := cegrpc.NewCryptoWalletClient(conn)
		defer closeResourcesFn()

		// CODE UNDER TEST
		response, err := performGRPCRequest(ctx, client, "Login", req)

		// EXPECTATION
		require.NoError(t, err)
		require.NotNil(t, response)
		require.NotNil(t, response.LoginResponse.Token)
		require.Equal(t, "Login successful", response.LoginResponse.Message)
		require.True(t, response.LoginResponse.Success)

		userMock.AssertExpectations(t)
	})

}
