package http_server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/glu-project/configs"
	"github.com/glu-project/idl/pb"
	"github.com/glu-project/utils/authenticate"
	"github.com/hashicorp/golang-lru/v2/expirable"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

type HttpServer struct {
	mux      *runtime.ServeMux
	server   *http.Server
	endpoint configs.Endpoint
}

func NewHttpServer(
	handler func(mux *runtime.ServeMux),
	endpoint configs.Endpoint,
	authenticator authenticate.Authenticator,
	client pb.UserServiceClient,
) *HttpServer {
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions:   protojson.MarshalOptions{UseEnumNumbers: false, EmitUnpopulated: true},
			UnmarshalOptions: protojson.UnmarshalOptions{AllowPartial: true},
		}),

		// runtime.WithErrorHandler(forwardErrorResponse),
	)
	fmt.Println("NewHttpServer")
	handler(mux)
	middlewares := []middlewareFunc{
		allowCORS,
		authorized(
			authenticator,
			expirable.NewLRU[string, []string](5, nil, 24*time.Hour),
		),
		httpLogger,
	}

	var handleR http.Handler = mux
	for _, handle := range middlewares {
		handleR = handle(handleR)
	}

	return &HttpServer{
		mux:      mux,
		endpoint: endpoint,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", endpoint.Port),
			Handler: handleR,
		},
	}
}

func (s *HttpServer) Start(ctx context.Context) error {
	log.Printf("Server listin in port: %d\n", s.endpoint.Port)
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *HttpServer) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
