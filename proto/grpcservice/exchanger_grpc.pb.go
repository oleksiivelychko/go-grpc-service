// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: proto/exchanger.proto

package grpcservice

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Exchanger_MakeExchange_FullMethodName = "/grpcservice.Exchanger/MakeExchange"
	Exchanger_Subscriber_FullMethodName   = "/grpcservice.Exchanger/Subscriber"
)

// ExchangerClient is the client API for Exchanger service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExchangerClient interface {
	MakeExchange(ctx context.Context, in *ExchangeRequest, opts ...grpc.CallOption) (*ExchangeResponse, error)
	// allows a client to subscribe for changes when the rate changes a response will be sent
	Subscriber(ctx context.Context, opts ...grpc.CallOption) (Exchanger_SubscriberClient, error)
}

type exchangerClient struct {
	cc grpc.ClientConnInterface
}

func NewExchangerClient(cc grpc.ClientConnInterface) ExchangerClient {
	return &exchangerClient{cc}
}

func (c *exchangerClient) MakeExchange(ctx context.Context, in *ExchangeRequest, opts ...grpc.CallOption) (*ExchangeResponse, error) {
	out := new(ExchangeResponse)
	err := c.cc.Invoke(ctx, Exchanger_MakeExchange_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exchangerClient) Subscriber(ctx context.Context, opts ...grpc.CallOption) (Exchanger_SubscriberClient, error) {
	stream, err := c.cc.NewStream(ctx, &Exchanger_ServiceDesc.Streams[0], Exchanger_Subscriber_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &exchangerSubscriberClient{stream}
	return x, nil
}

type Exchanger_SubscriberClient interface {
	Send(*ExchangeRequest) error
	Recv() (*StreamExchangeResponse, error)
	grpc.ClientStream
}

type exchangerSubscriberClient struct {
	grpc.ClientStream
}

func (x *exchangerSubscriberClient) Send(m *ExchangeRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *exchangerSubscriberClient) Recv() (*StreamExchangeResponse, error) {
	m := new(StreamExchangeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ExchangerServer is the server API for Exchanger service.
// All implementations should embed UnimplementedExchangerServer
// for forward compatibility
type ExchangerServer interface {
	MakeExchange(context.Context, *ExchangeRequest) (*ExchangeResponse, error)
	// allows a client to subscribe for changes when the rate changes a response will be sent
	Subscriber(Exchanger_SubscriberServer) error
}

// UnimplementedExchangerServer should be embedded to have forward compatible implementations.
type UnimplementedExchangerServer struct {
}

func (UnimplementedExchangerServer) MakeExchange(context.Context, *ExchangeRequest) (*ExchangeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MakeExchange not implemented")
}
func (UnimplementedExchangerServer) Subscriber(Exchanger_SubscriberServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscriber not implemented")
}

// UnsafeExchangerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExchangerServer will
// result in compilation errors.
type UnsafeExchangerServer interface {
	mustEmbedUnimplementedExchangerServer()
}

func RegisterExchangerServer(s grpc.ServiceRegistrar, srv ExchangerServer) {
	s.RegisterService(&Exchanger_ServiceDesc, srv)
}

func _Exchanger_MakeExchange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExchangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExchangerServer).MakeExchange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Exchanger_MakeExchange_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExchangerServer).MakeExchange(ctx, req.(*ExchangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Exchanger_Subscriber_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ExchangerServer).Subscriber(&exchangerSubscriberServer{stream})
}

type Exchanger_SubscriberServer interface {
	Send(*StreamExchangeResponse) error
	Recv() (*ExchangeRequest, error)
	grpc.ServerStream
}

type exchangerSubscriberServer struct {
	grpc.ServerStream
}

func (x *exchangerSubscriberServer) Send(m *StreamExchangeResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *exchangerSubscriberServer) Recv() (*ExchangeRequest, error) {
	m := new(ExchangeRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Exchanger_ServiceDesc is the grpc.ServiceDesc for Exchanger service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Exchanger_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpcservice.Exchanger",
	HandlerType: (*ExchangerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MakeExchange",
			Handler:    _Exchanger_MakeExchange_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscriber",
			Handler:       _Exchanger_Subscriber_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/exchanger.proto",
}