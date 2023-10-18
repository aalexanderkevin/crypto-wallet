package usecase

import (
	"context"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/container"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/repository"
	"github.com/aalexanderkevin/crypto-wallet/service"
)

type Webhook struct {
	config config.Config
	service.Bitcoin
	repository.Transaction
}

func NewWebhook(c *container.Container) *Webhook {
	return &Webhook{
		config:      c.Config(),
		Bitcoin:     c.Bitcoin(),
		Transaction: c.TransactionBtcRepo(),
	}
}

func (w Webhook) UpsertBitcoinTransaction(ctx context.Context, trx *model.Transaction) (err error) {
	logger := helper.GetLogger(ctx).WithField("method", "Usecase.Webhook.UpsertBitcoinTransaction")

	if _, err = w.Transaction.Upsert(ctx, trx); err != nil {
		logger.WithError(err).Warn("Failed upsert")
		return err
	}

	return nil
}
