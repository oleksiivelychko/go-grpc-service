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

func main() {
	logger := hclog.Default()
	gServer := grpc.NewServer()
	reflection.Register(gServer)
	sServer := server.NewServer(logger)

	gService.RegisterProductServer(gServer, sServer)
	gService.RegisterCurrencyServer(gServer, sServer)

	listen, err := net.Listen("tcp", "localhost:9091")
	if err != nil {
		logger.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	_ = gServer.Serve(listen)
}
