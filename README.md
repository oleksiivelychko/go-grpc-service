# go-grpc-protobuf

### gRPC service implementation is sharing the Protocol Buffer messages.

Install gRPCurl as `brew install grpcurl` and then test all in one:
```
grpcurl --plaintext localhost:9090 list
grpcurl --plaintext localhost:9090 list grpc_service.Product
grpcurl --plaintext localhost:9090 describe grpc_service.Product.GetProduct
grpcurl --plaintext localhost:9090 describe .grpc_service.ProductRequest
grpcurl --plaintext -d '{"id":1}' localhost:9090 grpc_service.Product.GetProduct
```

🎥 Thanks <a href="https://www.youtube.com/c/NicJackson">Nic Jackson</a> for sharing his knowledge.
