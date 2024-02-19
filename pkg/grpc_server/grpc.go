package grpc_server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/glu-project/configs"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	endpoint configs.Endpoint
	Server   *grpc.Server
}

func NewGrpcServer(endpoint configs.Endpoint) *GrpcServer {
	srv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_validator.StreamServerInterceptor(),
		)),
	)
	return &GrpcServer{
		endpoint: endpoint,
		Server:   srv,
	}
}

func (s *GrpcServer) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.endpoint.Port))
	if err != nil {
		return err
	}
	log.Printf("Server listening in port: %d\n", s.endpoint.Port)
	if err := s.Server.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (s *GrpcServer) Stop(ctx context.Context) error {
	s.Server.GracefulStop()
	return nil
}
