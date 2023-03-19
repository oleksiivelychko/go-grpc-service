package main

import (
	"github.com/oleksiivelychko/go-grpc-service/currency_server"
	"github.com/oleksiivelychko/go-grpc-service/exchange_processor"
	"github.com/oleksiivelychko/go-grpc-service/extractor_xml"
	"github.com/oleksiivelychko/go-grpc-service/proto/grpc_service"
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

	extractorXML := extractor_xml.NewExtractorXML(extractor_xml.SourceLocal, "rates.xml")
	exchangeProcessor, err := exchange_processor.NewExchangeProcessor(extractorXML)
	if err != nil {
		hcLogger.Error("unable to create exchangeProcessor", "error", err)
		os.Exit(1)
	}

	currencyServer := currency_server.NewCurrencyServer(hcLogger, exchangeProcessor)
	grpc_service.RegisterCurrencyServer(grpcServer, currencyServer)

	listenerTCP, err := net.Listen("tcp", localAddr)
	if err != nil {
		hcLogger.Error("unable to listen TCP", "error", err)
		os.Exit(1)
	}

	hcLogger.Info("starting gRPC server", "listening", localAddr)
	if err = grpcServer.Serve(listenerTCP); err != nil {
		hcLogger.Error("unable to start gRPC server", "error", err)
		os.Exit(1)
	}
}
