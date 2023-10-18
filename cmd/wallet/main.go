package main

import (
	"context"
	"os"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/container"
	"github.com/aalexanderkevin/crypto-wallet/repository/gormrepo"
	"github.com/aalexanderkevin/crypto-wallet/service"
	"github.com/aalexanderkevin/crypto-wallet/service/btc"
	"github.com/aalexanderkevin/crypto-wallet/service/eth"
	"github.com/aalexanderkevin/crypto-wallet/service/redis"
	"github.com/aalexanderkevin/crypto-wallet/service/trx"
	"github.com/aalexanderkevin/crypto-wallet/storage"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var rootCmd = &cobra.Command{
	Use:   "crypto-wallet",
	Short: "Crypto Wallet",
}

func init() {
	loadConfig()
	initLogging()
}

func main() {
	rootCmd := registerCommands(&defaultAppProvider{})
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}
}

func loadConfig() {
	err := config.Load()
	if err != nil {
		logrus.Errorf("Config error: %s", err.Error())
		os.Exit(1)
	}
}

func initLogging() *logrus.Logger {
	cfg := config.Instance()
	log := logrus.StandardLogger()
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	if strings.ToLower(cfg.LogFormat) == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	}
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)
	return log
}

func registerCommands(appProvider AppProvider) *cobra.Command {
	rootCmd.AddCommand(grpc(appProvider))
	rootCmd.AddCommand(restapi(appProvider))
	rootCmd.AddCommand(migrate(appProvider))

	return rootCmd
}

type AppProvider interface {
	// BuildContainer will return app dependencies container and resource clean up function
	// will return nil, nil, error when error happen
	BuildContainer(ctx context.Context, options buildOptions) (*container.Container, func(), error)
}

type buildOptions struct {
	Ethereum bool
	Bitcoin  bool
	Tron     bool
	Postgres bool
	Redis    bool
}

type defaultAppProvider struct {
}

func (defaultAppProvider) BuildContainer(ctx context.Context, options buildOptions) (*container.Container, func(), error) {
	var ethSvc service.Ethereum
	var btcSvc service.Bitcoin
	var trxSvc service.Tron
	var db *gorm.DB
	var redisSvc service.Cache

	cfg := config.Instance()

	// Init app container
	appContainer := container.NewContainer()
	appContainer.SetConfig(cfg)

	// Init Postgres
	if options.Postgres {
		db = storage.GetPostgresDb()
		appContainer.SetDb(db)

		userRepo := gormrepo.NewUserRepository(db)
		appContainer.SetUserRepo(userRepo)

		walletRepo := gormrepo.NewWalletRepository(db)
		appContainer.SetWalletRepo(walletRepo)

		transactionBtcRepo := gormrepo.NewBtcTransactionRepository(db)
		appContainer.SetTransactionBtcRepo(transactionBtcRepo)
		transactionTrxRepo := gormrepo.NewTrxTransactionRepository(db)
		appContainer.SetTransactionTrxRepo(transactionTrxRepo)
		transactionEthRepo := gormrepo.NewEthTransactionRepository(db)
		appContainer.SetTransactionEthRepo(transactionEthRepo)
	}

	// Init Service
	if options.Ethereum {
		ethSvc = eth.NewEthereumImpl(cfg)
		appContainer.SetEthereum(ethSvc)
	}

	if options.Bitcoin {
		btcSvc = btc.NewBitcoinImpl(cfg)
		appContainer.SetBitcoin(btcSvc)
	}

	if options.Tron {
		trxSvc = trx.NewTronImpl(cfg)
		appContainer.SetTron(trxSvc)
	}

	if options.Redis {
		redisSvc = redis.NewRedis(cfg.Redis)
		appContainer.SetRedis(redisSvc)
	}

	deferFn := func() {
		if ethSvc != nil {
			ethSvc.Close()
		}

		if db != nil {
			storage.CloseDB(db)
		}

		if options.Redis {
			redisSvc.Close()
		}

	}

	return appContainer, deferFn, nil
}
