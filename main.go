package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-grpc-protobuf/data"
	"github.com/oleksiivelychko/go-grpc-protobuf/processor"
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

	extractor := data.NewExtractorXml(data.SourceLocal)
	exchanger, err := processor.NewExchanger(extractor)
	if err != nil {
		logger.Error("unable to create exchanger", "error", err)
		os.Exit(1)
	}

	cServer := server.NewCurrencyServer(logger, exchanger)
	gService.RegisterCurrencyServer(gServer, cServer)

	listen, err := net.Listen("tcp", localAddr)
	if err != nil {
		logger.Error("unable to listen tcp", "error", err)
		os.Exit(1)
	}

	logger.Info("starting gRPC server", "listening", localAddr)
	err = gServer.Serve(listen)
	if err != nil {
		logger.Error("unable to start gRPC server", "error", err)
		os.Exit(1)
	}
}
