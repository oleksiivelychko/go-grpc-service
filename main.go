package main

import (
	"github.com/hashicorp/go-hclog"
	gService "github.com/oleksiivelychko/go-grpc-protobuf/proto/grpc_service"
	"github.com/oleksiivelychko/go-grpc-protobuf/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

const localAddr = "localhost:9091"

func main() {
	logger := hclog.Default()
	gServer := grpc.NewServer()
	reflection.Register(gServer)
	sServer := server.NewServer(logger)

	gService.RegisterProductServer(gServer, sServer)
	gService.RegisterCurrencyServer(gServer, sServer)

	listen, err := net.Listen("tcp", localAddr)
	if err != nil {
		logger.Error("unable to listen", "error", err)
		os.Exit(1)
	}

	logger.Info("starting gRPC server on", "addr", localAddr)
	err = gServer.Serve(listen)
	if err != nil {
		logger.Error("unable to start gRPC server", "error", err)
		os.Exit(1)
	}
}
