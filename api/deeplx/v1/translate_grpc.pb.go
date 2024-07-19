// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: deeplx/v1/translate.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	DeepLXService_Translate_FullMethodName = "/deeplx.v1.DeepLXService/Translate"
)

// DeepLXServiceClient is the client API for DeepLXService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeepLXServiceClient interface {
	Translate(ctx context.Context, in *TranslateRequest, opts ...grpc.CallOption) (*TranslationResult, error)
}

type deepLXServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDeepLXServiceClient(cc grpc.ClientConnInterface) DeepLXServiceClient {
	return &deepLXServiceClient{cc}
}

func (c *deepLXServiceClient) Translate(ctx context.Context, in *TranslateRequest, opts ...grpc.CallOption) (*TranslationResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TranslationResult)
	err := c.cc.Invoke(ctx, DeepLXService_Translate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeepLXServiceServer is the server API for DeepLXService service.
// All implementations must embed UnimplementedDeepLXServiceServer
// for forward compatibility
type DeepLXServiceServer interface {
	Translate(context.Context, *TranslateRequest) (*TranslationResult, error)
	mustEmbedUnimplementedDeepLXServiceServer()
}

// UnimplementedDeepLXServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDeepLXServiceServer struct {
}

func (UnimplementedDeepLXServiceServer) Translate(context.Context, *TranslateRequest) (*TranslationResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Translate not implemented")
}
func (UnimplementedDeepLXServiceServer) mustEmbedUnimplementedDeepLXServiceServer() {}

// UnsafeDeepLXServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeepLXServiceServer will
// result in compilation errors.
type UnsafeDeepLXServiceServer interface {
	mustEmbedUnimplementedDeepLXServiceServer()
}

func RegisterDeepLXServiceServer(s grpc.ServiceRegistrar, srv DeepLXServiceServer) {
	s.RegisterService(&DeepLXService_ServiceDesc, srv)
}

func _DeepLXService_Translate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TranslateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeepLXServiceServer).Translate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeepLXService_Translate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeepLXServiceServer).Translate(ctx, req.(*TranslateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DeepLXService_ServiceDesc is the grpc.ServiceDesc for DeepLXService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DeepLXService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "deeplx.v1.DeepLXService",
	HandlerType: (*DeepLXServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Translate",
			Handler:    _DeepLXService_Translate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "deeplx/v1/translate.proto",
}
