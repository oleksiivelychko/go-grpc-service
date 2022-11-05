package server

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-grpc-protobuf/processor"
	gService "github.com/oleksiivelychko/go-grpc-protobuf/proto/grpc_service"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"time"
)

type CurrencyServer struct {
	log           hclog.Logger
	exchanger     *processor.Exchanger
	subscriptions map[gService.Currency_SubscriberServer][]*gService.ExchangeRequest
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
		cs.log.Info("processing updated rates...")
		// loop over the subscribed clients
		for k, v := range cs.subscriptions {
			// loop over the subscribed rates
			for _, exchangeRequest := range v {
				fromCurrency := exchangeRequest.GetFrom().String()
				toCurrency := exchangeRequest.GetTo().String()

				rate, err := cs.exchanger.GetRate(fromCurrency, toCurrency)
				if err != nil {
					cs.log.Error("unable to get update of rate", "from", fromCurrency, "to", toCurrency)
				}

				err = k.Send(&gService.ExchangeResponse{From: exchangeRequest.From, To: exchangeRequest.To, Rate: rate})
				if err != nil {
					cs.log.Error("unable to send update of rate", "from", fromCurrency, "to", toCurrency, "rate", rate)
				}
			}
		}
	}
}

func (cs *CurrencyServer) MakeExchange(_ context.Context, r *gService.ExchangeRequest) (*gService.ExchangeResponse, error) {
	cs.log.Info("handle `grpc_service.Currency.MakeExchange`", "from", r.GetFrom(), "to", r.GetTo())

	rate, err := cs.exchanger.GetRate(r.GetFrom().String(), r.GetTo().String())
	if err != nil {
		cs.log.Error("cannot get rate", "error", err)
	}

	return &gService.ExchangeResponse{
		Rate:      rate,
		From:      r.GetFrom(),
		To:        r.GetTo(),
		UpdatedAt: timestamppb.Now(),
	}, nil
}

/*
Subscriber handles client messages.
*/
func (cs *CurrencyServer) Subscriber(srv gService.Currency_SubscriberServer) error {
	for {
		exchangeRequest, err := srv.Recv() // is blocking method
		if err == io.EOF {
			cs.log.Error("client has closed connection")
			break
		}

		if err != nil {
			cs.log.Error("unable to read from client", "error", err)
			return err
		}

		cs.log.Info("handle client request", "From", exchangeRequest.GetFrom(), "To", exchangeRequest.GetTo())

		subscriptionRequests, ok := cs.subscriptions[srv]
		if !ok {
			subscriptionRequests = []*gService.ExchangeRequest{}
		}

		subscriptionRequests = append(subscriptionRequests, exchangeRequest)
		cs.subscriptions[srv] = subscriptionRequests
	}

	return nil
}

/*
*
Requires to compile without 'require_unimplemented_servers'
*/
func (cs *CurrencyServer) mustEmbedUnimplementedCurrencyServer() {
	cs.log.Info("implement mustEmbedUnimplementedCurrencyServer for backward compatibility")
}
