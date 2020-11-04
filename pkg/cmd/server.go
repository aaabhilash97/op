package cmd

import (
	"context"
	"fmt"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/aaabhilash97/op/pkg/config"
	"github.com/aaabhilash97/op/pkg/protocol/grpc"
	"github.com/aaabhilash97/op/pkg/protocol/rest"
	v1 "github.com/aaabhilash97/op/pkg/service/v1"
)

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration

	if len(config.Server.GrpcPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", config.Server.GrpcPort)
	}

	if len(config.Server.HttpPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", config.Server.HttpPort)
	}

	v1API := v1.NewToDoServiceServer()

	// run HTTP gateway
	go func() {
		_ = rest.RunServer(ctx, config.Server.GrpcPort, config.Server.HttpPort)
	}()

	return grpc.RunServer(ctx, v1API, config.Server.GrpcPort)
}
