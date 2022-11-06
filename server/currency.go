package server

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-grpc-protobuf/processor"
	gService "github.com/oleksiivelychko/go-grpc-protobuf/proto/grpc_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

type CurrencyServer struct {
	log               hclog.Logger
	exchanger         *processor.Exchanger
	subscribedClients map[gService.Currency_SubscriberServer][]*gService.ExchangeRequest
}

func NewCurrencyServer(l hclog.Logger, e *processor.Exchanger) *CurrencyServer {
	requests := make(map[gService.Currency_SubscriberServer][]*gService.ExchangeRequest)
	cServer := &CurrencyServer{l, e, requests}

	go cServer.handleUpdates()

	return cServer
}

func (cs *CurrencyServer) handleUpdates() {
	updates := cs.exchanger.TrackRates(5 * time.Second)
	for range updates {
		cs.log.Info("handling updates...")
		// loop over the subscribed clients
		for subscribedClient, exchangeRequests := range cs.subscribedClients {
			// loop over the subscribed rates
			for _, exchangeRequest := range exchangeRequests {
				fromCurrency := exchangeRequest.GetFrom().String()
				toCurrency := exchangeRequest.GetTo().String()

				rate, err := cs.exchanger.GetRate(fromCurrency, toCurrency)
				if err != nil {
					cs.log.Error("unable to get update of rate", "from", fromCurrency, "to", toCurrency)
				}

				err = subscribedClient.Send(&gService.ExchangeResponse{
					From:      exchangeRequest.From,
					To:        exchangeRequest.To,
					Rate:      rate,
					CreatedAt: cs.exchanger.GetProtoTime(),
				})

				if err != nil {
					cs.log.Error("unable to send updated rate", "from", fromCurrency, "to", toCurrency, "rate", rate)
				}
			}
		}
	}
}

func (cs *CurrencyServer) MakeExchange(_ context.Context, r *gService.ExchangeRequest) (*gService.ExchangeResponse, error) {
	cs.log.Info("handle `grpc_service.Currency.MakeExchange`", "from", r.GetFrom(), "to", r.GetTo())

	if r.GetFrom() == r.GetTo() {
		grpcErr := status.Newf(
			codes.InvalidArgument,
			"base currency '%s' cannot be the same as destination '%s'",
			r.GetFrom(),
			r.GetTo(),
		)

		grpcStatus, err := grpcErr.WithDetails(r)
		if err != nil {
			return nil, err
		}

		return nil, grpcStatus.Err()
	}

	rate, err := cs.exchanger.GetRate(r.GetFrom().String(), r.GetTo().String())
	if err != nil {
		cs.log.Error("cannot get rate", "error", err)
	}

	return &gService.ExchangeResponse{
		Rate:      rate,
		From:      r.GetFrom(),
		To:        r.GetTo(),
		CreatedAt: cs.exchanger.GetProtoTime(),
	}, nil
}

/*
Subscriber implements the gRPC bidirectional streaming method.
*/
func (cs *CurrencyServer) Subscriber(srv gService.Currency_SubscriberServer) error {
	// handle client messages
	for {
		// 'Recv' is a blocking method which returns on client data.
		exchangeRequest, err := srv.Recv()
		if err == io.EOF {
			cs.log.Error("client has closed the connection")
			break
		}

		if err != nil {
			cs.log.Error("unable to read from client", "error", err)
			return err
		}

		cs.log.Info("handle client request", "From", exchangeRequest.GetFrom(), "To", exchangeRequest.GetTo())

		subscribedClient, ok := cs.subscribedClients[srv]
		if !ok {
			subscribedClient = []*gService.ExchangeRequest{}
		}

		subscribedClient = append(subscribedClient, exchangeRequest)
		cs.subscribedClients[srv] = subscribedClient
	}

	return nil
}

/*
*
mustEmbedUnimplementedCurrencyServer is required to compile without 'require_unimplemented_servers'.
*/
func (cs *CurrencyServer) mustEmbedUnimplementedCurrencyServer() {
	cs.log.Info("implement mustEmbedUnimplementedCurrencyServer for backward compatibility")
}
