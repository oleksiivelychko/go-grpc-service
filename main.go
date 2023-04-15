package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-grpc-service/exchanger"
	"github.com/oleksiivelychko/go-grpc-service/extractor"
	"github.com/oleksiivelychko/go-grpc-service/proto/grpcservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

func main() {
	hcLogger := hclog.New(&hclog.LoggerOptions{
		Name:       "go-grpc-service",
		Level:      hclog.LevelFromString("DEBUG"),
		Color:      1,
		TimeFormat: "02/01/2006 15:04:05",
	})

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	puller := extractor.New(extractor.SourceLocal, "rates.xml")
	processor, err := exchanger.NewProcessor(puller)
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
