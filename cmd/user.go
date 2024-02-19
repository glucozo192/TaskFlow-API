package cmd

import (
	"context"
	"log"
	"os"

	"github.com/glu-project/idl/pb"
	"github.com/glu-project/pkg/grpc_server"
	"github.com/spf13/cobra"
)

// userCmd represents the User command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		loadUser(ctx)

		errChan := make(chan error)
		start(ctx, errChan)
		go func() {
			select {
			case <-ctx.Done():
				log.Printf("Server initialization be canceled")
				os.Exit(1)
			case err := <-errChan:
				log.Printf("Server have an error: %v, server stop now\n", err.Error())
				if err := stop(ctx); err != nil {
					panic(err)
				}
			}
		}()
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func loadUser(ctx context.Context) {
	srv.loadPostgres(ctx)

	if err := srv.loadGrpcClients(ctx); err != nil {
		panic(err)
	}

	if err := srv.loadUserServices(ctx); err != nil {
		panic(err)
	}

	srv.userServer = grpc_server.NewGrpcServer(srv.cfg.UserServiceEndpoint)
	pb.RegisterUserServiceServer(srv.userServer.Server, srv.userService)

	// srv.userServer = grpc_server.NewGrpcServer(srv.cfg.UserlServiceEndpoint)

	srv.processors = append(srv.processors, srv.userServer)
	srv.factories = append(srv.factories, srv.postgresClient)
}
