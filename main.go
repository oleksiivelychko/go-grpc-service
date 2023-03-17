package main

import (
	cs "github.com/oleksiivelychko/go-grpc-service/currency_server"
	ep "github.com/oleksiivelychko/go-grpc-service/exchange_processor"
	gs "github.com/oleksiivelychko/go-grpc-service/proto/grpc_service"
	xe "github.com/oleksiivelychko/go-grpc-service/xml_extractor"
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

	xmlExtractor := xe.NewXmlExtractor(xe.SourceLocal)
	exchanger, err := ep.NewExchangeProcessor(xmlExtractor)
	if err != nil {
		hcLogger.Error("unable to create exchanger", "error", err)
		os.Exit(1)
	}

	currencyServer := cs.NewCurrencyServer(hcLogger, exchanger)
	gs.RegisterCurrencyServer(grpcServer, currencyServer)

	listen, err := net.Listen("tcp", localAddr)
	if err != nil {
		hcLogger.Error("unable to listen tcp", "error", err)
		os.Exit(1)
	}

	hcLogger.Info("starting gRPC cs", "listening", localAddr)
	if err = grpcServer.Serve(listen); err != nil {
		hcLogger.Error("unable to start gRPC cs", "error", err)
		os.Exit(1)
	}
}
