package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/storage"

	"github.com/spf13/cobra"
)

var (
	migrationPath  string
	rollback       bool
	versionToForce int
)

func migrate(appProvider AppProvider) *cobra.Command {
	cliCommand := &cobra.Command{
		Use:   "migrate",
		Short: "Run the database migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			app, closeResourcesFn, err := appProvider.BuildContainer(ctx, buildOptions{
				Postgres: true,
			})
			if err != nil {
				return err
			}
			if closeResourcesFn != nil {
				defer closeResourcesFn()
			}

			dB, _ := app.Db().DB()
			err = storage.MigratePostgresDb(dB, &migrationPath, rollback, versionToForce)
			if err != nil {
				fmt.Printf("Error when migration: %s \n", err.Error())
				os.Exit(1)
			}

			fmt.Println("Finish migrating database")
			return nil
		},
	}

	cfg := config.Instance()
	cliCommand.Flags().StringVarP(&migrationPath, "path", "p", cfg.Postgres.Migration.Path, "The migration folder")
	cliCommand.Flags().BoolVarP(&rollback, "rollback", "r", false, "Rollback to prev migration (-1 step)")
	cliCommand.Flags().IntVarP(&versionToForce, "force", "f", -1, "Force to specific version")
	return cliCommand
}
