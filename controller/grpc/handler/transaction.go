package handler

import (
	"context"
	"errors"

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

type Transaction struct {
	appContainer *container.Container
}

func NewTransactionHandler(appContainer *container.Container) *Transaction {
	return &Transaction{appContainer: appContainer}
}

func (w *Transaction) SendToken(ctx context.Context, r *cegrpc.SendRequest) (*cegrpc.SendResponse, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Handler.Transaction.SendToken")

	email := middleware.GetJWTData(ctx)
	if email == "" {
		err := errors.New("cant find id on token")
		logger.WithError(err)
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	req := &model.SendToken{
		Email:           helper.Pointer(email),
		ReceiverAddress: helper.Pointer(r.GetToAddress()),
		Amount:          helper.Pointer(r.GetAmount()),
		Token:           helper.Pointer(r.GetToken()),
	}
	err := req.Validate()
	if err != nil {
		logger.WithError(err).Warning("missing required field")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	transactionUseCase := usecase.NewTransaction(w.appContainer)
	var hashTx *string

	switch *req.Token {
	case "btc", "bitcoin":
		hashTx, err = transactionUseCase.SendBitcoin(ctx, req)
		if err != nil {
			return nil, response.SendErrorResponse(err)
		}
	case "trx", "tron":
		hashTx, err = transactionUseCase.SendTron(ctx, req)
		if err != nil {
			return nil, response.SendErrorResponse(err)
		}
	case "eth", "ethereum":
		hashTx, err = transactionUseCase.SendEthereum(ctx, req)
		if err != nil {
			return nil, response.SendErrorResponse(err)
		}
	default:
		err = errors.New("invalid transfer token")
		return nil, response.SendErrorResponse(err)
	}

	// Successful authentication, return hash transaction
	return &cegrpc.SendResponse{
		HashTransaction: *hashTx,
	}, nil
}
