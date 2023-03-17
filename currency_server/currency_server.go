package currency_server

import (
	"context"
	"github.com/hashicorp/go-hclog"
	ep "github.com/oleksiivelychko/go-grpc-service/exchange_processor"
	gs "github.com/oleksiivelychko/go-grpc-service/proto/grpc_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

type CurrencyServer struct {
	logger            hclog.Logger
	exchangeProcessor *ep.ExchangeProcessor
	subscribedClients map[gs.Currency_SubscriberServer][]*gs.ExchangeRequest
}

func NewCurrencyServer(logger hclog.Logger, exchanger *ep.ExchangeProcessor) *CurrencyServer {
	exchangeRequests := make(map[gs.Currency_SubscriberServer][]*gs.ExchangeRequest)

	currencyServer := &CurrencyServer{logger, exchanger, exchangeRequests}

	go currencyServer.handleUpdates()
	return currencyServer
}

func (currencyServer *CurrencyServer) handleUpdates() {
	updates := currencyServer.exchangeProcessor.TrackRates(5 * time.Second)
	for range updates {
		currencyServer.logger.Info("handling updates...")

		for subscribedClient, exchangeRequests := range currencyServer.subscribedClients {
			for _, exchangeRequest := range exchangeRequests {
				fromCurrency := exchangeRequest.GetFrom().String()
				toCurrency := exchangeRequest.GetTo().String()

				rate, err := currencyServer.exchangeProcessor.GetRate(fromCurrency, toCurrency)
				if err != nil {
					currencyServer.logger.Error("unable to get update of rate", "from", fromCurrency, "to", toCurrency)
				}

				err = subscribedClient.Send(&gs.StreamExchangeResponse{
					Message: &gs.StreamExchangeResponse_ExchangeResponse{
						ExchangeResponse: &gs.ExchangeResponse{
							From:      exchangeRequest.GetFrom(),
							To:        exchangeRequest.GetTo(),
							Rate:      rate,
							CreatedAt: currencyServer.exchangeProcessor.GetProtoTimestamp(),
						},
					},
				})

				if err != nil {
					currencyServer.logger.Error("unable to send updated rate", "from", fromCurrency, "to", toCurrency, "rate", rate)
				}
			}
		}
	}
}

func (currencyServer *CurrencyServer) MakeExchange(_ context.Context, exchangeRequest *gs.ExchangeRequest) (*gs.ExchangeResponse, error) {
	currencyServer.logger.Info("handle `grpc_service.Currency.MakeExchange`", "from", exchangeRequest.GetFrom(), "to", exchangeRequest.GetTo())

	if exchangeRequest.GetFrom() == exchangeRequest.GetTo() {
		grpcErr := status.Newf(
			codes.InvalidArgument,
			"base currency '%s' cannot be the same as destination '%s'",
			exchangeRequest.GetFrom(),
			exchangeRequest.GetTo(),
		)

		grpcStatus, err := grpcErr.WithDetails(exchangeRequest)
		if err != nil {
			return nil, err
		}

		return nil, grpcStatus.Err()
	}

	rate, err := currencyServer.exchangeProcessor.GetRate(exchangeRequest.GetFrom().String(), exchangeRequest.GetTo().String())
	if err != nil {
		currencyServer.logger.Error("cannot get rate", "error", err)
	}

	return &gs.ExchangeResponse{
		Rate:      rate,
		From:      exchangeRequest.GetFrom(),
		To:        exchangeRequest.GetTo(),
		CreatedAt: currencyServer.exchangeProcessor.GetProtoTimestamp(),
	}, nil
}

/*
Subscriber implements the gRPC bidirectional streaming method.
*/
func (currencyServer *CurrencyServer) Subscriber(subscriberServer gs.Currency_SubscriberServer) error {
	// handle client messages
	for {
		// 'Recv' is a blocking method which returns on client data.
		exchangeRequest, err := subscriberServer.Recv()
		if err == io.EOF {
			currencyServer.logger.Error("client has closed the connection")
			break
		}

		if err != nil {
			currencyServer.logger.Error("unable to read from client", "error", err)
			return err
		}

		currencyServer.logger.Info("handle client request", "From", exchangeRequest.GetFrom(), "To", exchangeRequest.GetTo())

		subscribedClient, ok := currencyServer.subscribedClients[subscriberServer]
		if !ok {
			subscribedClient = []*gs.ExchangeRequest{}
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
					currencyServer.logger.Error("unable to original request as metadata", "error", err)
				}

				break
			}
		}

		if validationErr != nil {
			err = subscriberServer.Send(&gs.StreamExchangeResponse{
				Message: &gs.StreamExchangeResponse_Error{
					Error: validationErr.Proto(),
				},
			})

			if err != nil {
				currencyServer.logger.Error("unable to send validation message", "error", err)
			}

			continue
		}

		subscribedClient = append(subscribedClient, exchangeRequest)
		currencyServer.subscribedClients[subscriberServer] = subscribedClient
	}

	return nil
}

/*
*
mustEmbedUnimplementedCurrencyServer is required to compile without 'require_unimplemented_servers'.
*/
func (currencyServer *CurrencyServer) mustEmbedUnimplementedCurrencyServer() {
	currencyServer.logger.Info("implement mustEmbedUnimplementedCurrencyServer for backward compatibility")
}
