package handler

import (
	"context"

	"github.com/aalexanderkevin/crypto-wallet/container"
	"github.com/aalexanderkevin/crypto-wallet/controller/grpc/response"
	"github.com/aalexanderkevin/crypto-wallet/controller/middleware"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/model"
	cegrpc "github.com/aalexanderkevin/crypto-wallet/transport/grpc/crypto-wallet"
	"github.com/aalexanderkevin/crypto-wallet/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	appContainer *container.Container
}

func NewUserHandler(appContainer *container.Container) *User {
	return &User{appContainer: appContainer}
}

func (u *User) Register(ctx context.Context, r *cegrpc.RegisterRequest) (*cegrpc.RegisterResponse, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Handler.User.Login")

	req := &model.User{
		Username: helper.Pointer(r.GetUsername()),
		Email:    helper.Pointer(r.GetEmail()),
		FullName: helper.Pointer(r.GetFullname()),
		Password: helper.Pointer(r.GetPassword()),
	}

	if err := req.Validate(); err != nil {
		logger.WithError(err).Warning("missing required field")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userUseCase := usecase.NewUser(u.appContainer)
	_, err := userUseCase.Register(ctx, req)
	if err != nil {
		return nil, response.SendErrorResponse(err)
	}

	// Successful authentication, generate token and return response
	return &cegrpc.RegisterResponse{
		Success: true,
		Message: "Register successful",
	}, nil
}

func (u *User) Login(ctx context.Context, r *cegrpc.LoginRequest) (*cegrpc.LoginResponse, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Handler.User.Login")
	config := u.appContainer.Config()

	req := &model.User{
		Username: helper.Pointer(r.GetUsername()),
		Password: helper.Pointer(r.GetPassword()),
	}

	if err := req.ValidateLogin(); err != nil {
		logger.WithError(err).Warning("missing required field")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userUseCase := usecase.NewUser(u.appContainer)
	user, err := userUseCase.Login(ctx, req)
	if err != nil {
		return nil, response.SendErrorResponse(err)
	}

	token, err := middleware.GenerateJwt(*user.Email, config.JwtSecret)
	if err != nil {
		logger.WithError(err).Warning("error generate JWT")
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	// Successful authentication, generate token and return response
	return &cegrpc.LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   *token,
	}, nil
}
