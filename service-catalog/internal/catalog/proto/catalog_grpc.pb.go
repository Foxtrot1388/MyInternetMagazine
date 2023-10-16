// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: proto/catalog.proto

package catalog

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

// CatalogApiClient is the client API for CatalogApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CatalogApiClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	Delete(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	List(ctx context.Context, in *ListParams, opts ...grpc.CallOption) (*ListResponse, error)
	Ping(ctx context.Context, in *PingParams, opts ...grpc.CallOption) (*PingResponse, error)
}

type catalogApiClient struct {
	cc grpc.ClientConnInterface
}

func NewCatalogApiClient(cc grpc.ClientConnInterface) CatalogApiClient {
	return &catalogApiClient{cc}
}

func (c *catalogApiClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/catalog.CatalogApi/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogApiClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/catalog.CatalogApi/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogApiClient) Delete(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, "/catalog.CatalogApi/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogApiClient) List(ctx context.Context, in *ListParams, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/catalog.CatalogApi/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogApiClient) Ping(ctx context.Context, in *PingParams, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/catalog.CatalogApi/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CatalogApiServer is the server API for CatalogApi service.
// All implementations must embed UnimplementedCatalogApiServer
// for forward compatibility
type CatalogApiServer interface {
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	Delete(context.Context, *GetRequest) (*DeleteResponse, error)
	List(context.Context, *ListParams) (*ListResponse, error)
	Ping(context.Context, *PingParams) (*PingResponse, error)
	mustEmbedUnimplementedCatalogApiServer()
}

// UnimplementedCatalogApiServer must be embedded to have forward compatible implementations.
type UnimplementedCatalogApiServer struct {
}

func (UnimplementedCatalogApiServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedCatalogApiServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedCatalogApiServer) Delete(context.Context, *GetRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedCatalogApiServer) List(context.Context, *ListParams) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedCatalogApiServer) Ping(context.Context, *PingParams) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedCatalogApiServer) mustEmbedUnimplementedCatalogApiServer() {}

// UnsafeCatalogApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CatalogApiServer will
// result in compilation errors.
type UnsafeCatalogApiServer interface {
	mustEmbedUnimplementedCatalogApiServer()
}

func RegisterCatalogApiServer(s grpc.ServiceRegistrar, srv CatalogApiServer) {
	s.RegisterService(&CatalogApi_ServiceDesc, srv)
}

func _CatalogApi_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogApiServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/catalog.CatalogApi/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogApiServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatalogApi_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogApiServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/catalog.CatalogApi/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogApiServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatalogApi_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogApiServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/catalog.CatalogApi/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogApiServer).Delete(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatalogApi_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogApiServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/catalog.CatalogApi/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogApiServer).List(ctx, req.(*ListParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatalogApi_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogApiServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/catalog.CatalogApi/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogApiServer).Ping(ctx, req.(*PingParams))
	}
	return interceptor(ctx, in, info, handler)
}

// CatalogApi_ServiceDesc is the grpc.ServiceDesc for CatalogApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CatalogApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "catalog.CatalogApi",
	HandlerType: (*CatalogApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _CatalogApi_Get_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _CatalogApi_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _CatalogApi_Delete_Handler,
		},
		{
			MethodName: "List",
			Handler:    _CatalogApi_List_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _CatalogApi_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/catalog.proto",
}
