package main

import (
	"github.com/oleksiivelychko/go-grpc-service/exchanger"
	"github.com/oleksiivelychko/go-grpc-service/extractor"
	"github.com/oleksiivelychko/go-grpc-service/logger"
	"github.com/oleksiivelychko/go-grpc-service/proto/grpcservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

func main() {
	log := logger.New()

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	extractorXML := extractor.New(extractor.SourceLocal, "rates.xml")
	processor, err := exchanger.NewProcessor(extractorXML)
	if err != nil {
		log.Error("unable to create processor: %s", err)
		os.Exit(1)
	}

	exchangerServer := exchanger.NewServer(processor, log)
	grpcservice.RegisterExchangerServer(grpcServer, exchangerServer)

	listenerTCP, err := net.Listen("tcp", exchangerServer.EnvAddress())
	if err != nil {
		log.Error("unable to listen TCP: %s", err)
		os.Exit(1)
	}

	log.Info("starting gRPC server on %s", exchangerServer.EnvAddress())
	if err = grpcServer.Serve(listenerTCP); err != nil {
		log.Error("unable to start gRPC server: %s", err)
		os.Exit(1)
	}
}
