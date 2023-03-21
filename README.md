# go-grpc-service

### The gRPC service serves external requests are sharing the Protocol Buffer messages.

📌 Install gRPCurl client and Protobuf compiler before use:
```
brew install grpcurl
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

📌 After compiling of **proto** files use **gRPCurl** to test it as:
```
grpcurl --plaintext localhost:9091 list
grpcurl --plaintext localhost:9091 list grpc_service.Currency
grpcurl --plaintext localhost:9091 describe grpc_service.Currency.MakeExchange
grpcurl --plaintext localhost:9091 describe .grpc_service.ExchangeRequest
grpcurl --plaintext -d '{"From": "EUR", "To": "USD"}' localhost:9091 grpc_service.Currency.MakeExchange
grpcurl --plaintext --msg-template -d @ localhost:9091 describe .grpc_service.ExchangeRequest
grpcurl --plaintext --msg-template -d @ localhost:9091 grpc_service.Currency.Subscriber
```

💡 Template message (might be inserted into stream as is):
```
{
  "From": "EUR",
  "To": "USD"
}
```

🎥 Thanks [Nic Jackson](https://www.youtube.com/c/NicJackson) for sharing his knowledge.
