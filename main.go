package main

import (
	"github.com/oleksiivelychko/go-grpc-service/exchanger"
	"github.com/oleksiivelychko/go-grpc-service/extractor"
	"github.com/oleksiivelychko/go-grpc-service/proto/grpcservice"
	"github.com/oleksiivelychko/go-utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

func main() {
	hcLogger := logger.NewHashicorp("go-grpc-service")
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	pullerXML := extractor.NewPullerXML(extractor.SourceLocal, "rates.xml")
	processor, err := exchanger.NewProcessor(pullerXML)
	if err != nil {
		hcLogger.Error("unable to create processor", "error", err)
		os.Exit(1)
	}

	exchangerServer := exchanger.NewServer(hcLogger, processor)
	grpcservice.RegisterExchangerServer(grpcServer, exchangerServer)

	listenerTCP, err := net.Listen("tcp", exchangerServer.EnvAddress())
	if err != nil {
		hcLogger.Error("unable to listen TCP", "error", err)
		os.Exit(1)
	}

	hcLogger.Info("starting gRPC server", "listening", exchangerServer.EnvAddress())
	if err = grpcServer.Serve(listenerTCP); err != nil {
		hcLogger.Error("unable to start gRPC server", "error", err)
		os.Exit(1)
	}
}
