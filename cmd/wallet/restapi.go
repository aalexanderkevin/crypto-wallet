package main

import (
	"context"

	controllerrestapi "github.com/aalexanderkevin/crypto-wallet/controller/restapi"
	"github.com/aalexanderkevin/crypto-wallet/helper"

	"github.com/segmentio/ksuid"
	"github.com/spf13/cobra"
)

func restapi(appProvider AppProvider) *cobra.Command {
	cliCommand := &cobra.Command{
		Use:   "run-restapi",
		Short: "Run REST API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := helper.ContextWithRequestId(context.Background(), ksuid.New().String())
			logger := helper.GetLogger(ctx).WithField("method", "server")

			app, closeResourcesFn, err := appProvider.BuildContainer(ctx, buildOptions{
				Postgres: true,
				Ethereum: true,
				Bitcoin:  true,
			})
			if err != nil {
				panic(err)
			}
			if closeResourcesFn != nil {
				defer closeResourcesFn()
			}

			// Start Http Server
			err = controllerrestapi.NewHttpServer(app).Start()
			if err != nil {
				logger.WithError(err).Error("Error starting web server")
				return err
			}

			return nil
		},
	}
	return cliCommand
}
