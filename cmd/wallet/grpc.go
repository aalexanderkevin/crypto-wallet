package main

import (
	"context"

	"github.com/aalexanderkevin/crypto-wallet/config"
	controllergrpc "github.com/aalexanderkevin/crypto-wallet/controller/grpc"

	"github.com/spf13/cobra"
)

func grpc(appProvider AppProvider) *cobra.Command {
	cliCommand := &cobra.Command{
		Use:   "run-grpc",
		Short: "Run gRPC server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			cfg := config.Instance()

			app, closeResourcesFn, err := appProvider.BuildContainer(ctx, buildOptions{
				Postgres: true,
				Redis:    true,
				Ethereum: true,
				Bitcoin:  true,
				Tron:     true,
			})
			if err != nil {
				return err
			}
			if closeResourcesFn != nil {
				defer closeResourcesFn()
			}

			controllergrpc.StartgRPC(app, cfg)
			return nil
		},
	}
	return cliCommand
}
