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
)

type Watcher struct {
	appContainer *container.Container
}

func NewWatcherHandler(appContainer *container.Container) *Watcher {
	return &Watcher{appContainer: appContainer}
}

func (w *Watcher) TriggerWatcher(ctx context.Context, req *cegrpc.TriggerWatcherRequest) (*cegrpc.TriggerWatcherResponse, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Handler.Watcher.TriggerWatcher")

	email := middleware.GetJWTData(ctx)
	if email == "" {
		err := errors.New("cant find id on token")
		logger.WithError(err)
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	transactionUseCase := usecase.NewTransaction(w.appContainer)
	watcherUseCase := usecase.NewWatcher(w.appContainer, *transactionUseCase)

	var address *string
	var err error
	switch req.GetToken() {
	case "eth", "ethereum":
		address, err = watcherUseCase.TriggerWatcherEth(ctx, &email)
		if err != nil {
			return nil, response.SendErrorResponse(err)
		}
	case "trx", "tron":
		address, err = watcherUseCase.TriggerWatcherTrx(ctx, &email)
		if err != nil {
			return nil, response.SendErrorResponse(err)
		}
	default:
		err := errors.New("invalid transfer token")
		return nil, response.SendErrorResponse(err)
	}

	// Successful authentication, return hash transaction
	return &cegrpc.TriggerWatcherResponse{
		Address: *address,
	}, nil
}
