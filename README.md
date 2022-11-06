# go-grpc-service

### The gRPC service serves external requests are sharing the Protocol Buffer messages.

Install gRPCurl as `brew install grpcurl` and then test it as:
```
grpcurl --plaintext localhost:9091 list
grpcurl --plaintext localhost:9091 list grpc_service.Currency
grpcurl --plaintext localhost:9091 describe grpc_service.Currency.MakeExchange
grpcurl --plaintext localhost:9091 describe .grpc_service.ExchangeRequest
grpcurl --plaintext -d '{"From": "EUR", "To": "USD"}' localhost:9091 grpc_service.Currency.MakeExchange
grpcurl --plaintext --msg-template -d @ localhost:9091 describe .grpc_service.ExchangeRequest
grpcurl --plaintext --msg-template -d @ localhost:9091 grpc_service.Currency.Subscriber
```

Template message (might be inserted into stream as is):
```
{
  "From": "EUR",
  "To": "USD"
}
```

🎥 Thanks <a href="https://www.youtube.com/c/NicJackson">Nic Jackson</a> for sharing his knowledge.
