package server

import (
	"context"
	"github.com/hashicorp/go-hclog"
	gService "github.com/oleksiivelychko/go-grpc-protobuf/proto/grpc_service"
)

type Server struct {
	log hclog.Logger
}

func NewServer(l hclog.Logger) *Server {
	return &Server{l}
}

/*
GetProduct returns the sample product.
*/
func (s *Server) GetProduct(_ context.Context, pr *gService.ProductRequest) (*gService.ProductResponse, error) {
	s.log.Info("handle 'grpc_service.Product.GetProduct'", "ID", pr.GetId())

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
func (s *Server) mustEmbedUnimplementedProductServer() {
	s.log.Info("implement mustEmbedUnimplementedProductServer for backward compatibility")
}
