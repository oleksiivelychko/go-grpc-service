# go-grpc-protobuf

### gRPC service implementation is sharing the Protocol Buffer messages.

Install gRPCurl as `brew install grpcurl` and then test all in one:
```
grpcurl --plaintext localhost:9091 list
grpcurl --plaintext localhost:9091 list grpc_service.Currency
grpcurl --plaintext localhost:9091 describe grpc_service.Currency.MakeExchange
grpcurl --plaintext localhost:9091 describe .grpc_service.ExchangeRequest
grpcurl --plaintext -d '{"from": "EUR", "to": "USD"}' localhost:9091 grpc_service.Currency.MakeExchange
grpcurl --plaintext --msg-template -d @ localhost:9091 describe .grpc_service.ExchangeRequest
grpcurl --plaintext --msg-template -d @ localhost:9091 grpc_service.Currency.Subscriber
```

Template message:
```
{
  "From": "EUR",
  "To": "USD"
}
```

🎥 Thanks <a href="https://www.youtube.com/c/NicJackson">Nic Jackson</a> for sharing his knowledge.
