// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: profile.proto

package profile

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

// ProfileApiClient is the client API for ProfileApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProfileApiClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	Delete(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	Ping(ctx context.Context, in *PingParams, opts ...grpc.CallOption) (*PingResponse, error)
}

type profileApiClient struct {
	cc grpc.ClientConnInterface
}

func NewProfileApiClient(cc grpc.ClientConnInterface) ProfileApiClient {
	return &profileApiClient{cc}
}

func (c *profileApiClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/profile.ProfileApi/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileApiClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/profile.ProfileApi/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileApiClient) Delete(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, "/profile.ProfileApi/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileApiClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/profile.ProfileApi/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileApiClient) Ping(ctx context.Context, in *PingParams, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/profile.ProfileApi/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfileApiServer is the server API for ProfileApi service.
// All implementations must embed UnimplementedProfileApiServer
// for forward compatibility
type ProfileApiServer interface {
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	Delete(context.Context, *GetRequest) (*DeleteResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	Ping(context.Context, *PingParams) (*PingResponse, error)
	mustEmbedUnimplementedProfileApiServer()
}

// UnimplementedProfileApiServer must be embedded to have forward compatible implementations.
type UnimplementedProfileApiServer struct {
}

func (UnimplementedProfileApiServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedProfileApiServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedProfileApiServer) Delete(context.Context, *GetRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedProfileApiServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedProfileApiServer) Ping(context.Context, *PingParams) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedProfileApiServer) mustEmbedUnimplementedProfileApiServer() {}

// UnsafeProfileApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProfileApiServer will
// result in compilation errors.
type UnsafeProfileApiServer interface {
	mustEmbedUnimplementedProfileApiServer()
}

func RegisterProfileApiServer(s grpc.ServiceRegistrar, srv ProfileApiServer) {
	s.RegisterService(&ProfileApi_ServiceDesc, srv)
}

func _ProfileApi_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileApiServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.ProfileApi/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileApiServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileApi_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileApiServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.ProfileApi/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileApiServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileApi_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileApiServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.ProfileApi/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileApiServer).Delete(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileApi_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileApiServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.ProfileApi/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileApiServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileApi_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileApiServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.ProfileApi/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileApiServer).Ping(ctx, req.(*PingParams))
	}
	return interceptor(ctx, in, info, handler)
}

// ProfileApi_ServiceDesc is the grpc.ServiceDesc for ProfileApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProfileApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "profile.ProfileApi",
	HandlerType: (*ProfileApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _ProfileApi_Get_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _ProfileApi_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ProfileApi_Delete_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _ProfileApi_Login_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _ProfileApi_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "profile.proto",
}
