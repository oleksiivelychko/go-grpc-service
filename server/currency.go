package server

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-grpc-service/processor"
	gService "github.com/oleksiivelychko/go-grpc-service/proto/grpc_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

type CurrencyServer struct {
	logger            hclog.Logger
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
		cs.logger.Info("handling updates...")

		for subscribedClient, exchangeRequests := range cs.subscribedClients {
			for _, exchangeRequest := range exchangeRequests {
				fromCurrency := exchangeRequest.GetFrom().String()
				toCurrency := exchangeRequest.GetTo().String()

				rate, err := cs.exchanger.GetRate(fromCurrency, toCurrency)
				if err != nil {
					cs.logger.Error("unable to get update of rate", "from", fromCurrency, "to", toCurrency)
				}

				err = subscribedClient.Send(&gService.StreamExchangeResponse{
					Message: &gService.StreamExchangeResponse_ExchangeResponse{
						ExchangeResponse: &gService.ExchangeResponse{
							From:      exchangeRequest.GetFrom(),
							To:        exchangeRequest.GetTo(),
							Rate:      rate,
							CreatedAt: cs.exchanger.GetProtoTime(),
						},
					},
				})

				if err != nil {
					cs.logger.Error("unable to send updated rate", "from", fromCurrency, "to", toCurrency, "rate", rate)
				}
			}
		}
	}
}

func (cs *CurrencyServer) MakeExchange(_ context.Context, r *gService.ExchangeRequest) (*gService.ExchangeResponse, error) {
	cs.logger.Info("handle `grpc_service.Currency.MakeExchange`", "from", r.GetFrom(), "to", r.GetTo())

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
		cs.logger.Error("cannot get rate", "error", err)
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
			cs.logger.Error("client has closed the connection")
			break
		}

		if err != nil {
			cs.logger.Error("unable to read from client", "error", err)
			return err
		}

		cs.logger.Info("handle client request", "From", exchangeRequest.GetFrom(), "To", exchangeRequest.GetTo())

		subscribedClient, ok := cs.subscribedClients[srv]
		if !ok {
			subscribedClient = []*gService.ExchangeRequest{}
		}

		var validationErr *status.Status
		// check that subscriber does not exist
		for _, client := range subscribedClient {
			if client.From == exchangeRequest.From && client.To == exchangeRequest.To {
				// subscriber exists, return error
				validationErr = status.Newf(
					codes.AlreadyExists,
					"unable to subscribe for currency as subscription already exists: base '%s', destination '%s'",
					exchangeRequest.GetFrom(),
					exchangeRequest.GetTo(),
				)

				if validationErr, err = validationErr.WithDetails(exchangeRequest); err != nil {
					cs.logger.Error("unable to original request as metadata", "error", err)
				}

				break
			}
		}

		if validationErr != nil {
			err = srv.Send(&gService.StreamExchangeResponse{
				Message: &gService.StreamExchangeResponse_Error{
					Error: validationErr.Proto(),
				},
			})

			if err != nil {
				cs.logger.Error("unable to send validation message", "error", err)
			}

			continue
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
	cs.logger.Info("implement mustEmbedUnimplementedCurrencyServer for backward compatibility")
}
