package handler_test

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/container"
	grpccontroller "github.com/aalexanderkevin/crypto-wallet/controller/grpc"
	"github.com/aalexanderkevin/crypto-wallet/controller/grpc/handler"
	"github.com/aalexanderkevin/crypto-wallet/controller/middleware"
	cegrpc "github.com/aalexanderkevin/crypto-wallet/transport/grpc/crypto-wallet"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// CustomResponse encapsulates all possible response types.
type CustomResponse struct {
	LoginResponse    *cegrpc.LoginResponse
	RegisterResponse *cegrpc.RegisterResponse
	SendResponse     *cegrpc.SendResponse
}

func SetupGRPCConn(t *testing.T, ctx context.Context, cb func(appContainer *container.Container) *container.Container) (*grpc.ClientConn, func()) {
	cfg := config.Instance()
	appContainer := DefaultAppContainer()
	if cb != nil {
		appContainer = cb(appContainer)
	}

	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	// List of excluded methods (full method names).
	excludedMethods := []string{
		"/crypto_wallet.CryptoWallet/Login",    // Exclude Login method
		"/crypto_wallet.CryptoWallet/Register", // Exclude Register method
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.JWTMiddleware(cfg.JwtSecret, excludedMethods)),
	)

	controllers := &grpccontroller.Controllers{
		User:        *handler.NewUserHandler(appContainer),
		Wallet:      *handler.NewWalletHandler(appContainer),
		Transaction: *handler.NewTransactionHandler(appContainer),
		Watcher:     *handler.NewWatcherHandler(appContainer),
	}
	cegrpc.RegisterCryptoWalletServer(server, controllers)

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	// Create a context and gRPC client using the listener
	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	deferFn := func() {
		server.GracefulStop()
		lis.Close()
		conn.Close()
	}

	return conn, deferFn
}

func performGRPCRequest(ctx context.Context, client cegrpc.CryptoWalletClient, method string, requestMessage interface{}) (*CustomResponse, error) {
	switch method {
	case "Login":
		response, err := client.Login(ctx, requestMessage.(*cegrpc.LoginRequest))
		if err != nil {
			return nil, err
		}
		return &CustomResponse{LoginResponse: response}, nil
	case "Register":
		response, err := client.Register(ctx, requestMessage.(*cegrpc.RegisterRequest))
		if err != nil {
			return nil, err
		}
		return &CustomResponse{RegisterResponse: response}, nil
	case "SendToken":
		response, err := client.SendToken(ctx, requestMessage.(*cegrpc.SendRequest))
		if err != nil {
			return nil, err
		}
		return &CustomResponse{SendResponse: response}, nil
	default:
		return nil, errors.New("Unsupported gRPC method")
	}
}

func DefaultAppContainer() *container.Container {
	appContainer := container.NewContainer()
	appContainer.SetConfig(config.Instance())

	return appContainer
}
