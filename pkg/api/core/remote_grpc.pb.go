// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: remote.proto

package core

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
	RemoteService_Connect_FullMethodName      = "/remote.RemoteService/Connect"
	RemoteService_DisConnect_FullMethodName   = "/remote.RemoteService/DisConnect"
	RemoteService_Remote_FullMethodName       = "/remote.RemoteService/Remote"
	RemoteService_RemoteInput_FullMethodName  = "/remote.RemoteService/RemoteInput"
	RemoteService_RemoteOutput_FullMethodName = "/remote.RemoteService/RemoteOutput"
)

// RemoteServiceClient is the client API for RemoteService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RemoteServiceClient interface {
	Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*ConnectResponse, error)
	DisConnect(ctx context.Context, in *DisconnectRequest, opts ...grpc.CallOption) (*Result, error)
	Remote(ctx context.Context, opts ...grpc.CallOption) (RemoteService_RemoteClient, error)
	RemoteInput(ctx context.Context, in *RemoteRequest, opts ...grpc.CallOption) (*Result, error)
	RemoteOutput(ctx context.Context, in *RemoteOutputRequest, opts ...grpc.CallOption) (RemoteService_RemoteOutputClient, error)
}

type remoteServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRemoteServiceClient(cc grpc.ClientConnInterface) RemoteServiceClient {
	return &remoteServiceClient{cc}
}

func (c *remoteServiceClient) Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*ConnectResponse, error) {
	out := new(ConnectResponse)
	err := c.cc.Invoke(ctx, RemoteService_Connect_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *remoteServiceClient) DisConnect(ctx context.Context, in *DisconnectRequest, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, RemoteService_DisConnect_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *remoteServiceClient) Remote(ctx context.Context, opts ...grpc.CallOption) (RemoteService_RemoteClient, error) {
	stream, err := c.cc.NewStream(ctx, &RemoteService_ServiceDesc.Streams[0], RemoteService_Remote_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &remoteServiceRemoteClient{stream}
	return x, nil
}

type RemoteService_RemoteClient interface {
	Send(*RemoteRequest) error
	Recv() (*RemoteResponse, error)
	grpc.ClientStream
}

type remoteServiceRemoteClient struct {
	grpc.ClientStream
}

func (x *remoteServiceRemoteClient) Send(m *RemoteRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *remoteServiceRemoteClient) Recv() (*RemoteResponse, error) {
	m := new(RemoteResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *remoteServiceClient) RemoteInput(ctx context.Context, in *RemoteRequest, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, RemoteService_RemoteInput_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *remoteServiceClient) RemoteOutput(ctx context.Context, in *RemoteOutputRequest, opts ...grpc.CallOption) (RemoteService_RemoteOutputClient, error) {
	stream, err := c.cc.NewStream(ctx, &RemoteService_ServiceDesc.Streams[1], RemoteService_RemoteOutput_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &remoteServiceRemoteOutputClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RemoteService_RemoteOutputClient interface {
	Recv() (*RemoteResponse, error)
	grpc.ClientStream
}

type remoteServiceRemoteOutputClient struct {
	grpc.ClientStream
}

func (x *remoteServiceRemoteOutputClient) Recv() (*RemoteResponse, error) {
	m := new(RemoteResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RemoteServiceServer is the server API for RemoteService service.
// All implementations must embed UnimplementedRemoteServiceServer
// for forward compatibility
type RemoteServiceServer interface {
	Connect(context.Context, *ConnectRequest) (*ConnectResponse, error)
	DisConnect(context.Context, *DisconnectRequest) (*Result, error)
	Remote(RemoteService_RemoteServer) error
	RemoteInput(context.Context, *RemoteRequest) (*Result, error)
	RemoteOutput(*RemoteOutputRequest, RemoteService_RemoteOutputServer) error
	mustEmbedUnimplementedRemoteServiceServer()
}

// UnimplementedRemoteServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRemoteServiceServer struct {
}

func (UnimplementedRemoteServiceServer) Connect(context.Context, *ConnectRequest) (*ConnectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedRemoteServiceServer) DisConnect(context.Context, *DisconnectRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisConnect not implemented")
}
func (UnimplementedRemoteServiceServer) Remote(RemoteService_RemoteServer) error {
	return status.Errorf(codes.Unimplemented, "method Remote not implemented")
}
func (UnimplementedRemoteServiceServer) RemoteInput(context.Context, *RemoteRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoteInput not implemented")
}
func (UnimplementedRemoteServiceServer) RemoteOutput(*RemoteOutputRequest, RemoteService_RemoteOutputServer) error {
	return status.Errorf(codes.Unimplemented, "method RemoteOutput not implemented")
}
func (UnimplementedRemoteServiceServer) mustEmbedUnimplementedRemoteServiceServer() {}

// UnsafeRemoteServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RemoteServiceServer will
// result in compilation errors.
type UnsafeRemoteServiceServer interface {
	mustEmbedUnimplementedRemoteServiceServer()
}

func RegisterRemoteServiceServer(s grpc.ServiceRegistrar, srv RemoteServiceServer) {
	s.RegisterService(&RemoteService_ServiceDesc, srv)
}

func _RemoteService_Connect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoteServiceServer).Connect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RemoteService_Connect_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoteServiceServer).Connect(ctx, req.(*ConnectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RemoteService_DisConnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisconnectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoteServiceServer).DisConnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RemoteService_DisConnect_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoteServiceServer).DisConnect(ctx, req.(*DisconnectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RemoteService_Remote_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RemoteServiceServer).Remote(&remoteServiceRemoteServer{stream})
}

type RemoteService_RemoteServer interface {
	Send(*RemoteResponse) error
	Recv() (*RemoteRequest, error)
	grpc.ServerStream
}

type remoteServiceRemoteServer struct {
	grpc.ServerStream
}

func (x *remoteServiceRemoteServer) Send(m *RemoteResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *remoteServiceRemoteServer) Recv() (*RemoteRequest, error) {
	m := new(RemoteRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _RemoteService_RemoteInput_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoteServiceServer).RemoteInput(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RemoteService_RemoteInput_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoteServiceServer).RemoteInput(ctx, req.(*RemoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RemoteService_RemoteOutput_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RemoteOutputRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RemoteServiceServer).RemoteOutput(m, &remoteServiceRemoteOutputServer{stream})
}

type RemoteService_RemoteOutputServer interface {
	Send(*RemoteResponse) error
	grpc.ServerStream
}

type remoteServiceRemoteOutputServer struct {
	grpc.ServerStream
}

func (x *remoteServiceRemoteOutputServer) Send(m *RemoteResponse) error {
	return x.ServerStream.SendMsg(m)
}

// RemoteService_ServiceDesc is the grpc.ServiceDesc for RemoteService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RemoteService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "remote.RemoteService",
	HandlerType: (*RemoteServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Connect",
			Handler:    _RemoteService_Connect_Handler,
		},
		{
			MethodName: "DisConnect",
			Handler:    _RemoteService_DisConnect_Handler,
		},
		{
			MethodName: "RemoteInput",
			Handler:    _RemoteService_RemoteInput_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Remote",
			Handler:       _RemoteService_Remote_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "RemoteOutput",
			Handler:       _RemoteService_RemoteOutput_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "remote.proto",
}
