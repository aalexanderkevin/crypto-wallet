package response

import (
	"errors"

	"github.com/aalexanderkevin/crypto-wallet/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SendErrorResponse(err error) error {
	var modelError model.Error
	if errors.As(err, &modelError) {
		return status.Error(modelError.Code, modelError.Message)
	}

	return status.Error(codes.Internal, err.Error())
}
