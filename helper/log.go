package helper

import (
	"context"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
)

const (
	ContextKeyRequestId   Key = "requestId"
	ContextKeyJwtData     Key = "jwtData"
	ContextKeyTokenBearer Key = "tokenBearer"
)

type Key string

func (k Key) String() string {
	return string(k)
}

func GetLogger(c context.Context) *logrus.Entry {
	reqId := c.Value(ContextKeyRequestId)
	if reqId == nil {
		reqId = c.Value(string(ContextKeyRequestId))
		if reqId == nil {
			reqId = ksuid.New().String()
		}
	}
	logger := logrus.WithField(ContextKeyRequestId.String(), reqId)

	return logger
}

func ContextWithRequestId(ctx context.Context, requestId string) context.Context {
	defaultContext := context.TODO()
	if ctx != nil {
		defaultContext = ctx
	}
	return context.WithValue(defaultContext, ContextKeyRequestId, requestId)
}
