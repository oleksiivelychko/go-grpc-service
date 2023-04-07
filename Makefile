HOST := localhost
PORT := 9091
PROTO := -proto ./proto/exchanger.proto

.PHONY: proto
proto:
	protoc --go_out=proto --go-grpc_out=require_unimplemented_servers=false:proto proto/*.proto

run:
	HOST=$(HOST) PORT=$(PORT) go run main.go

grpcurl-list-services:
	grpcurl -plaintext $(HOST):$(PORT) list

grpcurl-list-service-methods:
	grpcurl -plaintext $(PROTO) $(HOST):$(PORT) list grpcservice.Exchanger

grpcurl-describe-service-method:
	grpcurl -plaintext $(PROTO) $(HOST):$(PORT) describe grpcservice.Exchanger.MakeExchange

grpcurl-describe-message:
	grpcurl -plaintext $(PROTO) -msg-template $(HOST):$(PORT) describe .grpcservice.ExchangeRequest

grpcurl-send-message:
	grpcurl -plaintext -d '{"From": "EUR", "To": "USD"}' $(PROTO) $(HOST):$(PORT) grpcservice.Exchanger.MakeExchange

grpcurl-message-to-stream:
	grpcurl -plaintext $(PROTO) -msg-template -d @ $(HOST):$(PORT) grpcservice.Exchanger.Subscriber
