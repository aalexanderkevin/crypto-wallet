package handler

import (
	"context"
	"errors"

	"github.com/aalexanderkevin/crypto-wallet/container"
	"github.com/aalexanderkevin/crypto-wallet/controller/grpc/response"
	"github.com/aalexanderkevin/crypto-wallet/controller/middleware"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	cegrpc "github.com/aalexanderkevin/crypto-wallet/transport/grpc/crypto-wallet"
	"github.com/aalexanderkevin/crypto-wallet/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Wallet struct {
	appContainer *container.Container

	cegrpc.UnimplementedCryptoWalletServer
}

func NewWalletHandler(appContainer *container.Container) *Wallet {
	return &Wallet{appContainer: appContainer}
}

func (w *Wallet) CreateWallet(ctx context.Context, r *emptypb.Empty) (*cegrpc.CreteWalletResponse, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Handler.Wallet.CreateWallet")

	email := middleware.GetJWTData(ctx)
	if email == "" {
		err := errors.New("cant find email on token")
		logger.WithError(err)
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	walletUseCase := usecase.NewWallet(w.appContainer)
	wallet, err := walletUseCase.CreateNewWallet(ctx, &email)
	if err != nil {
		return nil, response.SendErrorResponse(err)
	}

	// Successful authentication, return hash wallet
	return &cegrpc.CreteWalletResponse{
		Id:         *wallet.Id,
		Email:      email,
		BtcAddress: *wallet.BtcAddress,
		EthAddress: *wallet.EthAddress,
		TrxAddress: *wallet.TrxAddress,
	}, nil
}
