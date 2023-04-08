package exchanger

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-grpc-service/proto/grpcservice"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"os"
	"time"
)

type Server struct {
	logger            hclog.Logger
	processor         *Processor
	subscribedClients map[grpcservice.Exchanger_SubscriberServer][]*grpcservice.ExchangeRequest
}

func NewServer(logger hclog.Logger, processor *Processor) *Server {
	requests := make(map[grpcservice.Exchanger_SubscriberServer][]*grpcservice.ExchangeRequest)
	server := &Server{logger, processor, requests}

	go server.handleUpdates()
	return server
}

func (server *Server) handleUpdates() {
	updates := server.processor.TrackRates(5 * time.Second)
	for range updates {
		server.logger.Info("handling updates...")

		for subscribedClient, exchangeRequests := range server.subscribedClients {
			for _, exchangeRequest := range exchangeRequests {
				fromCurrency := exchangeRequest.GetFrom().String()
				toCurrency := exchangeRequest.GetTo().String()

				rate, err := server.processor.GetRate(fromCurrency, toCurrency)
				if err != nil {
					server.logger.Error(
						"unable to get update of rate",
						"from",
						fromCurrency,
						"to",
						toCurrency,
					)
				}

				err = subscribedClient.Send(&grpcservice.StreamExchangeResponse{
					Message: &grpcservice.StreamExchangeResponse_ExchangeResponse{
						ExchangeResponse: &grpcservice.ExchangeResponse{
							From:      exchangeRequest.GetFrom(),
							To:        exchangeRequest.GetTo(),
							Rate:      rate,
							CreatedAt: server.processor.GetProtoTimestamp(),
						},
					},
				})

				if err != nil {
					server.logger.Error(
						"unable to send updated rate",
						"from",
						fromCurrency,
						"to",
						toCurrency,
						"rate",
						rate,
					)
				}
			}
		}
	}
}

func (server *Server) MakeExchange(_ context.Context, exchangeRequest *grpcservice.ExchangeRequest) (*grpcservice.ExchangeResponse, error) {
	server.logger.Info(
		"handle 'grpcservice.Exchanger.MakeExchange'",
		"from",
		exchangeRequest.GetFrom(),
		"to",
		exchangeRequest.GetTo(),
	)

	if exchangeRequest.GetFrom() == exchangeRequest.GetTo() {
		grpcErr := status.Newf(
			codes.InvalidArgument,
			"base currency %s cannot be the same as destination %s",
			exchangeRequest.GetFrom(),
			exchangeRequest.GetTo(),
		)

		grpcStatus, err := grpcErr.WithDetails(exchangeRequest)
		if err != nil {
			return nil, err
		}

		return nil, grpcStatus.Err()
	}

	rate, err := server.processor.GetRate(
		exchangeRequest.GetFrom().String(),
		exchangeRequest.GetTo().String(),
	)

	if err != nil {
		server.logger.Error("unable to get rate", "error", err)
	}

	return &grpcservice.ExchangeResponse{
		Rate:      rate,
		From:      exchangeRequest.GetFrom(),
		To:        exchangeRequest.GetTo(),
		CreatedAt: server.processor.GetProtoTimestamp(),
	}, nil
}

/*
Subscriber implements the gRPC bidirectional streaming method.
*/
func (server *Server) Subscriber(subscriberServer grpcservice.Exchanger_SubscriberServer) error {
	// handle client messages
	for {
		// 'Recv' is a blocking method which returns on client data
		exchangeRequest, err := subscriberServer.Recv()
		if err == io.EOF {
			server.logger.Error("client has closed the connection")
			break
		}

		if err != nil {
			server.logger.Error("unable to read from client", "error", err)
			return err
		}

		server.logger.Info(
			"handle client request",
			"From",
			exchangeRequest.GetFrom(),
			"To",
			exchangeRequest.GetTo(),
		)

		subscribedClient, ok := server.subscribedClients[subscriberServer]
		if !ok {
			subscribedClient = []*grpcservice.ExchangeRequest{}
		}

		var validationErr *status.Status
		// check that subscriber does not exist
		for _, client := range subscribedClient {
			if client.From == exchangeRequest.From && client.To == exchangeRequest.To {
				// subscriber exists, return error
				validationErr = status.Newf(
					codes.AlreadyExists,
					"unable to subscribe for currency as subscription already exists: base %s, destination %s",
					exchangeRequest.GetFrom(),
					exchangeRequest.GetTo(),
				)

				if validationErr, err = validationErr.WithDetails(exchangeRequest); err != nil {
					server.logger.Error("unable to get original request as metadata", "error", err)
				}

				break
			}
		}

		if validationErr != nil {
			err = subscriberServer.Send(&grpcservice.StreamExchangeResponse{
				Message: &grpcservice.StreamExchangeResponse_Error{
					Error: validationErr.Proto(),
				},
			})

			if err != nil {
				server.logger.Error("unable to send validation message", "error", err)
			}

			continue
		}

		subscribedClient = append(subscribedClient, exchangeRequest)
		server.subscribedClients[subscriberServer] = subscribedClient
	}

	return nil
}

func (server *Server) EnvAddress() string {
	return fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
}

/*
mustEmbedUnimplementedExchangerServer is required to compile without 'require_unimplemented_servers'.
*/
func (server *Server) mustEmbedUnimplementedExchangerServer() {
	server.logger.Info("implement mustEmbedUnimplementedExchangerServer for backward compatibility")
}
