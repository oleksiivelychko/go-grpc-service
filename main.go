package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-grpc-service/data"
	"github.com/oleksiivelychko/go-grpc-service/processor"
	gService "github.com/oleksiivelychko/go-grpc-service/proto/grpc_service"
	"github.com/oleksiivelychko/go-grpc-service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

const localAddr = "localhost:9091"

func main() {
	logger := hclog.New(&hclog.LoggerOptions{TimeFormat: "2006/01/02 15:04:05", Color: 1})
	gServer := grpc.NewServer()
	reflection.Register(gServer)

	extractor := data.NewExtractor(data.SourceLocal)
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
	if err = gServer.Serve(listen); err != nil {
		logger.Error("unable to start gRPC server", "error", err)
		os.Exit(1)
	}
}
