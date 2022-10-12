install-protobuf:
	brew install protobuf
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

.PHONY: proto
proto:
	protoc --go_out=proto --go-grpc_out=proto proto/*.proto