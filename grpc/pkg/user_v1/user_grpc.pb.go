// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.29.0--rc1
// source: user.proto

package user_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UserV1Client is the client API for UserV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserV1Client interface {
	// CreateUser creates new user
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	// GetUser returns user
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
	// UpdateUser updates user
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// DeleteUser deletes existing user
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type userV1Client struct {
	cc grpc.ClientConnInterface
}

func NewUserV1Client(cc grpc.ClientConnInterface) UserV1Client {
	return &userV1Client{cc}
}

func (c *userV1Client) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, "/user_v1.UserV1/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userV1Client) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	out := new(GetUserResponse)
	err := c.cc.Invoke(ctx, "/user_v1.UserV1/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userV1Client) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/user_v1.UserV1/UpdateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userV1Client) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/user_v1.UserV1/DeleteUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserV1Server is the server API for UserV1 service.
// All implementations must embed UnimplementedUserV1Server
// for forward compatibility
type UserV1Server interface {
	// CreateUser creates new user
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	// GetUser returns user
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	// UpdateUser updates user
	UpdateUser(context.Context, *UpdateUserRequest) (*emptypb.Empty, error)
	// DeleteUser deletes existing user
	DeleteUser(context.Context, *DeleteUserRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedUserV1Server()
}

// UnimplementedUserV1Server must be embedded to have forward compatible implementations.
type UnimplementedUserV1Server struct {
}

func (UnimplementedUserV1Server) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedUserV1Server) GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedUserV1Server) UpdateUser(context.Context, *UpdateUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedUserV1Server) DeleteUser(context.Context, *DeleteUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedUserV1Server) mustEmbedUnimplementedUserV1Server() {}

// UnsafeUserV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserV1Server will
// result in compilation errors.
type UnsafeUserV1Server interface {
	mustEmbedUnimplementedUserV1Server()
}

func RegisterUserV1Server(s grpc.ServiceRegistrar, srv UserV1Server) {
	s.RegisterService(&UserV1_ServiceDesc, srv)
}

func _UserV1_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserV1Server).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_v1.UserV1/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserV1Server).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserV1_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserV1Server).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_v1.UserV1/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserV1Server).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserV1_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserV1Server).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_v1.UserV1/UpdateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserV1Server).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserV1_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserV1Server).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_v1.UserV1/DeleteUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserV1Server).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserV1_ServiceDesc is the grpc.ServiceDesc for UserV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user_v1.UserV1",
	HandlerType: (*UserV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _UserV1_CreateUser_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _UserV1_GetUser_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _UserV1_UpdateUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _UserV1_DeleteUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}
