package main

import (
	"github.com/oleksiivelychko/go-grpc-service/processor"
	grpcService "github.com/oleksiivelychko/go-grpc-service/proto/grpc_service"
	"github.com/oleksiivelychko/go-grpc-service/server"
	extractor "github.com/oleksiivelychko/go-grpc-service/xml_extractor"
	"github.com/oleksiivelychko/go-utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

const localAddr = "localhost:9091"

func main() {
	hcLogger := logger.NewLogger("go-grpc-service")
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	xmlExtractor := extractor.NewXmlExtractor(extractor.SourceLocal)
	exchanger, err := processor.NewExchanger(xmlExtractor)
	if err != nil {
		hcLogger.Error("unable to create exchanger", "error", err)
		os.Exit(1)
	}

	currencyServer := server.NewCurrencyServer(hcLogger, exchanger)
	grpcService.RegisterCurrencyServer(grpcServer, currencyServer)

	listen, err := net.Listen("tcp", localAddr)
	if err != nil {
		hcLogger.Error("unable to listen tcp", "error", err)
		os.Exit(1)
	}

	hcLogger.Info("starting gRPC server", "listening", localAddr)
	if err = grpcServer.Serve(listen); err != nil {
		hcLogger.Error("unable to start gRPC server", "error", err)
		os.Exit(1)
	}
}
