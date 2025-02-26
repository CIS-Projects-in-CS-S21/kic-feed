// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// HealthTrackingClient is the client API for HealthTracking service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HealthTrackingClient interface {
	//Given health data obtained upon user request, said health data is returned to user.
	GetHealthDataForUser(ctx context.Context, in *GetHealthDataForUserRequest, opts ...grpc.CallOption) (*GetHealthDataForUserResponse, error)
	//Health data requested to be added by user is added, and error is returned if appropriate.
	AddHealthDataForUser(ctx context.Context, in *AddHealthDataForUserRequest, opts ...grpc.CallOption) (*AddHealthDataForUserResponse, error)
	//Health data requested by user to be deleted is deleted and said deleted entries are returned to user.
	DeleteHealthDataForUser(ctx context.Context, in *DeleteHealthDataForUserRequest, opts ...grpc.CallOption) (*DeleteHealthDataForUserResponse, error)
	//Health data requested to be updated by user is updated, and error is returned if appropriate.
	UpdateHealthDataForDate(ctx context.Context, in *UpdateHealthDataForDateRequest, opts ...grpc.CallOption) (*UpdateHealthDataForDateResponse, error)
	// Given user ID, returns a mental health score for a user
	GetMentalHealthScoreForUser(ctx context.Context, in *GetMentalHealthScoreForUserRequest, opts ...grpc.CallOption) (*GetMentalHealthScoreForUserResponse, error)
	// Given a date and user ID, return health data log for a specific date
	GetHealthDataByDate(ctx context.Context, in *GetHealthDataByDateRequest, opts ...grpc.CallOption) (*GethealthDataByDateResponse, error)
}

type healthTrackingClient struct {
	cc grpc.ClientConnInterface
}

func NewHealthTrackingClient(cc grpc.ClientConnInterface) HealthTrackingClient {
	return &healthTrackingClient{cc}
}

func (c *healthTrackingClient) GetHealthDataForUser(ctx context.Context, in *GetHealthDataForUserRequest, opts ...grpc.CallOption) (*GetHealthDataForUserResponse, error) {
	out := new(GetHealthDataForUserResponse)
	err := c.cc.Invoke(ctx, "/kic.health.HealthTracking/GetHealthDataForUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *healthTrackingClient) AddHealthDataForUser(ctx context.Context, in *AddHealthDataForUserRequest, opts ...grpc.CallOption) (*AddHealthDataForUserResponse, error) {
	out := new(AddHealthDataForUserResponse)
	err := c.cc.Invoke(ctx, "/kic.health.HealthTracking/AddHealthDataForUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *healthTrackingClient) DeleteHealthDataForUser(ctx context.Context, in *DeleteHealthDataForUserRequest, opts ...grpc.CallOption) (*DeleteHealthDataForUserResponse, error) {
	out := new(DeleteHealthDataForUserResponse)
	err := c.cc.Invoke(ctx, "/kic.health.HealthTracking/DeleteHealthDataForUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *healthTrackingClient) UpdateHealthDataForDate(ctx context.Context, in *UpdateHealthDataForDateRequest, opts ...grpc.CallOption) (*UpdateHealthDataForDateResponse, error) {
	out := new(UpdateHealthDataForDateResponse)
	err := c.cc.Invoke(ctx, "/kic.health.HealthTracking/UpdateHealthDataForDate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *healthTrackingClient) GetMentalHealthScoreForUser(ctx context.Context, in *GetMentalHealthScoreForUserRequest, opts ...grpc.CallOption) (*GetMentalHealthScoreForUserResponse, error) {
	out := new(GetMentalHealthScoreForUserResponse)
	err := c.cc.Invoke(ctx, "/kic.health.HealthTracking/GetMentalHealthScoreForUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *healthTrackingClient) GetHealthDataByDate(ctx context.Context, in *GetHealthDataByDateRequest, opts ...grpc.CallOption) (*GethealthDataByDateResponse, error) {
	out := new(GethealthDataByDateResponse)
	err := c.cc.Invoke(ctx, "/kic.health.HealthTracking/GetHealthDataByDate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HealthTrackingServer is the server API for HealthTracking service.
// All implementations must embed UnimplementedHealthTrackingServer
// for forward compatibility
type HealthTrackingServer interface {
	//Given health data obtained upon user request, said health data is returned to user.
	GetHealthDataForUser(context.Context, *GetHealthDataForUserRequest) (*GetHealthDataForUserResponse, error)
	//Health data requested to be added by user is added, and error is returned if appropriate.
	AddHealthDataForUser(context.Context, *AddHealthDataForUserRequest) (*AddHealthDataForUserResponse, error)
	//Health data requested by user to be deleted is deleted and said deleted entries are returned to user.
	DeleteHealthDataForUser(context.Context, *DeleteHealthDataForUserRequest) (*DeleteHealthDataForUserResponse, error)
	//Health data requested to be updated by user is updated, and error is returned if appropriate.
	UpdateHealthDataForDate(context.Context, *UpdateHealthDataForDateRequest) (*UpdateHealthDataForDateResponse, error)
	// Given user ID, returns a mental health score for a user
	GetMentalHealthScoreForUser(context.Context, *GetMentalHealthScoreForUserRequest) (*GetMentalHealthScoreForUserResponse, error)
	// Given a date and user ID, return health data log for a specific date
	GetHealthDataByDate(context.Context, *GetHealthDataByDateRequest) (*GethealthDataByDateResponse, error)
	mustEmbedUnimplementedHealthTrackingServer()
}

// UnimplementedHealthTrackingServer must be embedded to have forward compatible implementations.
type UnimplementedHealthTrackingServer struct {
}

func (UnimplementedHealthTrackingServer) GetHealthDataForUser(context.Context, *GetHealthDataForUserRequest) (*GetHealthDataForUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHealthDataForUser not implemented")
}
func (UnimplementedHealthTrackingServer) AddHealthDataForUser(context.Context, *AddHealthDataForUserRequest) (*AddHealthDataForUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddHealthDataForUser not implemented")
}
func (UnimplementedHealthTrackingServer) DeleteHealthDataForUser(context.Context, *DeleteHealthDataForUserRequest) (*DeleteHealthDataForUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteHealthDataForUser not implemented")
}
func (UnimplementedHealthTrackingServer) UpdateHealthDataForDate(context.Context, *UpdateHealthDataForDateRequest) (*UpdateHealthDataForDateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateHealthDataForDate not implemented")
}
func (UnimplementedHealthTrackingServer) GetMentalHealthScoreForUser(context.Context, *GetMentalHealthScoreForUserRequest) (*GetMentalHealthScoreForUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMentalHealthScoreForUser not implemented")
}
func (UnimplementedHealthTrackingServer) GetHealthDataByDate(context.Context, *GetHealthDataByDateRequest) (*GethealthDataByDateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHealthDataByDate not implemented")
}
func (UnimplementedHealthTrackingServer) mustEmbedUnimplementedHealthTrackingServer() {}

// UnsafeHealthTrackingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HealthTrackingServer will
// result in compilation errors.
type UnsafeHealthTrackingServer interface {
	mustEmbedUnimplementedHealthTrackingServer()
}

func RegisterHealthTrackingServer(s grpc.ServiceRegistrar, srv HealthTrackingServer) {
	s.RegisterService(&_HealthTracking_serviceDesc, srv)
}

func _HealthTracking_GetHealthDataForUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetHealthDataForUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthTrackingServer).GetHealthDataForUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kic.health.HealthTracking/GetHealthDataForUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthTrackingServer).GetHealthDataForUser(ctx, req.(*GetHealthDataForUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HealthTracking_AddHealthDataForUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddHealthDataForUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthTrackingServer).AddHealthDataForUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kic.health.HealthTracking/AddHealthDataForUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthTrackingServer).AddHealthDataForUser(ctx, req.(*AddHealthDataForUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HealthTracking_DeleteHealthDataForUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteHealthDataForUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthTrackingServer).DeleteHealthDataForUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kic.health.HealthTracking/DeleteHealthDataForUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthTrackingServer).DeleteHealthDataForUser(ctx, req.(*DeleteHealthDataForUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HealthTracking_UpdateHealthDataForDate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateHealthDataForDateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthTrackingServer).UpdateHealthDataForDate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kic.health.HealthTracking/UpdateHealthDataForDate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthTrackingServer).UpdateHealthDataForDate(ctx, req.(*UpdateHealthDataForDateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HealthTracking_GetMentalHealthScoreForUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMentalHealthScoreForUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthTrackingServer).GetMentalHealthScoreForUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kic.health.HealthTracking/GetMentalHealthScoreForUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthTrackingServer).GetMentalHealthScoreForUser(ctx, req.(*GetMentalHealthScoreForUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HealthTracking_GetHealthDataByDate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetHealthDataByDateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthTrackingServer).GetHealthDataByDate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kic.health.HealthTracking/GetHealthDataByDate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthTrackingServer).GetHealthDataByDate(ctx, req.(*GetHealthDataByDateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _HealthTracking_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kic.health.HealthTracking",
	HandlerType: (*HealthTrackingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetHealthDataForUser",
			Handler:    _HealthTracking_GetHealthDataForUser_Handler,
		},
		{
			MethodName: "AddHealthDataForUser",
			Handler:    _HealthTracking_AddHealthDataForUser_Handler,
		},
		{
			MethodName: "DeleteHealthDataForUser",
			Handler:    _HealthTracking_DeleteHealthDataForUser_Handler,
		},
		{
			MethodName: "UpdateHealthDataForDate",
			Handler:    _HealthTracking_UpdateHealthDataForDate_Handler,
		},
		{
			MethodName: "GetMentalHealthScoreForUser",
			Handler:    _HealthTracking_GetMentalHealthScoreForUser_Handler,
		},
		{
			MethodName: "GetHealthDataByDate",
			Handler:    _HealthTracking_GetHealthDataByDate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/health.proto",
}
