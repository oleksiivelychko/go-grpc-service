.PHONY: proto
proto:
	protoc --go_out=proto --go-grpc_out=require_unimplemented_servers=false:proto proto/*.proto
