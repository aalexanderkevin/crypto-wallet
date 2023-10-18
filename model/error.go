package model

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
)

const (
	ErrorBadRequest          codes.Code = codes.InvalidArgument
	ErrorUnauthenticated     codes.Code = codes.Unauthenticated
	ErrorNotFound            codes.Code = codes.NotFound
	ErrorDuplicate           codes.Code = codes.AlreadyExists
	ErrorUnprocessableEntity codes.Code = codes.InvalidArgument
	ErrorInternalServer      codes.Code = codes.Internal
)

type Error struct {
	Message string
	Code    codes.Code
}

func NewError(message string, code codes.Code) Error {
	return Error{
		Message: message,
		Code:    code,
	}
}

func (e Error) Error() string {
	return e.Message
}

func NewParameterError(msg *string) Error {
	defaultMessage := "invalid parameter"
	if msg == nil {
		msg = &defaultMessage
	}
	return NewError(*msg, ErrorUnprocessableEntity)
}

func NewNotFoundError() Error {
	return NewError("resource not found", ErrorNotFound)
}

func NewDuplicateError() Error {
	return NewError("resource already exists", ErrorDuplicate)
}

func NewUnauthenticatedError() Error {
	return NewError("unauthorized access", ErrorUnauthenticated)
}

func NewInvalidPasswordError() Error {
	return NewError("invalid password", ErrorUnauthenticated)
}

func NewBadRequestError(msg *string) Error {
	defaultMessage := "bad request"
	if msg == nil {
		msg = &defaultMessage
	}
	return NewError(*msg, ErrorBadRequest)
}

func IsDuplicateError(e error) bool {
	var internalErr Error
	if !errors.As(e, &internalErr) {
		return false
	}

	return internalErr.Code == ErrorDuplicate
}

func IsNotFoundError(e error) bool {
	var internalErr Error
	if !errors.As(e, &internalErr) {
		return false
	}

	return internalErr.Code == ErrorNotFound
}

func NewStatusNotOKError(code int, body []byte) Error {
	e := fmt.Sprintf("status is not ok, status=%d body=%s", code, body)
	return NewError(e, ErrorInternalServer)
}

func IsParameterError(e error) bool {
	var internalErr Error
	if !errors.As(e, &internalErr) {
		return false
	}

	return internalErr.Code == ErrorUnprocessableEntity
}

func IsBadRequestError(e error) bool {
	var internalErr Error
	if !errors.As(e, &internalErr) {
		return false
	}

	return internalErr.Code == ErrorBadRequest
}
