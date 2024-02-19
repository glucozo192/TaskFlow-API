package cmd

import (
	"context"
	"log"

	"github.com/glu-project/idl/pb"
	"github.com/glu-project/pkg/http_server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
)

var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		loadGateway(ctx)

		errChan := make(chan error)
		start(ctx, errChan)
		err := <-errChan
		log.Printf("Server have an error: %v, server stop now\n", err.Error())
		if err := stop(ctx); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gatewayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gatewayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gatewayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func loadGateway(ctx context.Context) {
	if err := srv.loadGrpcClients(ctx); err != nil {
		panic(err)
	}
	srv.loadDefault(ctx)

	srv.gatewayServer = http_server.NewHttpServer(func(mux *runtime.ServeMux) {
		if err := pb.RegisterUserServiceHandlerClient(ctx, mux, srv.userClient); err != nil {
			panic(err)
		}

	}, srv.cfg.GatewayServiceEndpoint,
		srv.authenticator,
		srv.userClient)

	srv.processors = append(srv.processors, srv.gatewayServer)
	srv.factories = append(srv.factories, srv.userConnClient)
}
