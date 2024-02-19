package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/glu-project/idl/pb"
	userService "github.com/glu-project/internal/user/services"
	"github.com/glu-project/pkg/grpc_client"
	"github.com/glu-project/pkg/grpc_server"
	"github.com/glu-project/pkg/http_server"
	postgres_client "github.com/glu-project/pkg/postgres"
	"github.com/glu-project/utils/authenticate"

	"github.com/glu-project/configs"
	"github.com/glu-project/utils"
)

var srv server

type processor interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type factory interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}

type server struct {
	cfg *configs.Config

	processors []processor
	factories  []factory

	userServer    *grpc_server.GrpcServer
	gatewayServer *http_server.HttpServer

	authenticator authenticate.Authenticator
	// // repository
	// userRepo    userRepo.UserRepo
	// productRepo productRepo.ProductRepo

	// // services
	// productSrv product.ProductServiceServer
	userService pb.UserServiceServer

	// config
	config utils.Config

	// database clients
	postgresClient *postgres_client.PostgresClient
	//conn           *pgx.Conn

	// grpc clients
	userClient pb.UserServiceClient

	// grpc conn clients
	userConnClient *grpc_client.GrpcClient
}

// type processor interface {
// 	Init(ctx context.Context) error
// 	Start(ctx context.Context) error
// 	Stop(ctx context.Context) error
// }

// type factory interface {
// 	Connect(ctx context.Context) error
// 	Stop(ctx context.Context) error
// }

func (s *server) loadGrpcClients(ctx context.Context) error {
	//connections
	s.userConnClient = grpc_client.NewGrpcClient(s.cfg.UserServiceEndpoint)

	// clients for user service
	s.userClient = pb.NewUserServiceClient(s.userConnClient)
	return nil
}

func (s *server) loadUserServices(ctx context.Context) error {
	var err error
	s.authenticator, err = authenticate.NewPasetoAuthenticator(s.cfg.SymmetricKey, 24*time.Hour)
	if err != nil {
		panic(err)
	}
	s.userService = userService.NewUserService(
		s.postgresClient,
	)
	return nil
}

func (s *server) LoadConfig(ctx context.Context) error {
	var err error
	s.config, err = utils.LoadConfig(".")
	if err != nil {
		return err
	}
	return nil
}

func (s *server) loadDefault(ctx context.Context) {
	var err error
	s.authenticator, err = authenticate.NewPasetoAuthenticator(s.cfg.SymmetricKey, 24*time.Hour)
	if err != nil {
		panic(err)
	}
}

func (s *server) loadPostgres(ctx context.Context) {
	var err error
	srv.postgresClient, err = postgres_client.NewClient(srv.cfg.PostgresUrl)
	fmt.Println("connect to postgres:", srv.cfg.PostgresUrl)
	if err != nil {
		panic(err)
	}
}

func start(ctx context.Context, errChan chan error) {
	for _, f := range srv.factories {
		if err := f.Connect(ctx); err != nil {
			errChan <- err
		}
	}
	for _, p := range srv.processors {
		func(p processor) {
			if err := p.Start(ctx); err != nil {
				errChan <- err
			}
		}(p)
	}
}

func stop(ctx context.Context) error {
	for _, processor := range srv.processors {
		if err := processor.Stop(ctx); err != nil {
			return err
		}
	}

	for _, database := range srv.factories {
		if err := database.Close(ctx); err != nil {
			return err
		}
	}
	return nil
}
