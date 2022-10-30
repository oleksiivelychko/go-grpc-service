package server

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-grpc-protobuf/data"
	gService "github.com/oleksiivelychko/go-grpc-protobuf/proto/grpc_service"
)

type CurrencyServer struct {
	log hclog.Logger
	e   *data.Exchanger
}

func NewCurrencyServer(l hclog.Logger, e *data.Exchanger) *CurrencyServer {
	return &CurrencyServer{l, e}
}

func (cs *CurrencyServer) MakeExchange(_ context.Context, r *gService.ExchangeRequest) (*gService.ExchangeResponse, error) {
	cs.log.Info("handle `grpc_service.Currency.MakeExchange`", "from", r.GetFrom(), "to", r.GetTo())

	rate, err := cs.e.GetRate(r.GetFrom(), r.GetTo())
	if err != nil {
		cs.log.Error("cannot get rate", "error", err)
	}

	return &gService.ExchangeResponse{
		Rate: rate,
	}, nil
}
