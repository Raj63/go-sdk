package main

import (
	"log"

	sdkgrpc "github.com/Raj63/go-sdk/grpc"
	"github.com/Raj63/go-sdk/logger"
	"google.golang.org/grpc"
)

func main() {
	_logger := logger.NewLogger()

	// Setup the gRPC/HTTP server.
	server, err := sdkgrpc.NewServer(
		&sdkgrpc.ServerConfig{},
		_logger,
		func() error {
			return nil
		},
		[]grpc.UnaryServerInterceptor{}...,
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Register the gRPC server implementation.
	// api.RegisterCustomServiceServer(
	// 	server.GRPCServer(),
	// 	&pkggrpc.Server{},
	// )

	// finally serve the server
	if err := server.Serve(); err != nil {
		_logger.Fatal(err)
	}

}
