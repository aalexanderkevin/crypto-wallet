package grpc

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/container"
	"github.com/aalexanderkevin/crypto-wallet/controller/grpc/handler"
	"github.com/aalexanderkevin/crypto-wallet/controller/middleware"
	cegrpc "github.com/aalexanderkevin/crypto-wallet/transport/grpc/crypto-wallet"

	"google.golang.org/grpc"
)

type Controllers struct {
	handler.User
	handler.Wallet
	handler.Transaction
	handler.Watcher
}

func StartgRPC(app *container.Container, cfg config.Config) {
	// Start gRPC
	lis, err := net.Listen("tcp", cfg.Service.Host+":"+cfg.Service.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// List of excluded methods (full method names).
	excludedMethods := []string{
		"/crypto_wallet.CryptoWallet/Login",    // Exclude Login method
		"/crypto_wallet.CryptoWallet/Register", // Exclude Register method
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.JWTMiddleware(cfg.JwtSecret, excludedMethods)),
	)

	controllers := &Controllers{
		User:        *handler.NewUserHandler(app),
		Wallet:      *handler.NewWalletHandler(app),
		Transaction: *handler.NewTransactionHandler(app),
		Watcher:     *handler.NewWatcherHandler(app),
	}
	cegrpc.RegisterCryptoWalletServer(server, controllers)

	// Listen for OS signals to gracefully stop the server
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("gRPC server listening at " + lis.Addr().String())
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for a signal to gracefully stop the server
	<-stopCh

	// Stop the server gracefully
	log.Println("Shutting down gRPC server...")
	server.GracefulStop()
	log.Println("gRPC server has been gracefully stopped.")
}
