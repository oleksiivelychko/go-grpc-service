package server

import (
	"context"
	"github.com/hashicorp/go-hclog"
	gService "github.com/oleksiivelychko/go-grpc-protobuf/proto/grpc_service"
)

type ProductServer struct {
	log hclog.Logger
}

func NewProductServer(l hclog.Logger) *ProductServer {
	return &ProductServer{l}
}

/*
GetProduct returns the sample product.
*/
func (s *ProductServer) GetProduct(_ context.Context, r *gService.ProductRequest) (*gService.ProductResponse, error) {
	s.log.Info("handle `grpc_service.Product.GetProduct`", "id", r.GetId())

	department := gService.ProductResponse_Department{
		Number: "D1",
		Type:   gService.ProductResponse_MarketType(2),
	}

	return &gService.ProductResponse{
		Name:        "Main product",
		Sku:         "000-000-000",
		Price:       0.99,
		Departments: []*gService.ProductResponse_Department{&department},
	}, nil
}

/*
*
Requires to compile without 'require_unimplemented_servers'
*/
func (s *ProductServer) mustEmbedUnimplementedProductServer() {
	s.log.Info("implement mustEmbedUnimplementedProductServer for backward compatibility")
}
