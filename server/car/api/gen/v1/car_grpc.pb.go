// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.2
// source: car.proto

package carpb

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

// CarServiceClient is the client API for CarService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CarServiceClient interface {
	CreateCar(ctx context.Context, in *CreateCarRequest, opts ...grpc.CallOption) (*CarEntity, error)
	GetCar(ctx context.Context, in *GetCarRequest, opts ...grpc.CallOption) (*Car, error)
	GetCars(ctx context.Context, in *GetCarsRequest, opts ...grpc.CallOption) (*GetCarsResponse, error)
	LockCar(ctx context.Context, in *LockCarRequest, opts ...grpc.CallOption) (*LockCarResponse, error)
	UnLockCar(ctx context.Context, in *UnLockCarRequest, opts ...grpc.CallOption) (*UnLockCarResponse, error)
	UpdateCar(ctx context.Context, in *UpdateCarRequest, opts ...grpc.CallOption) (*UpdateCarResponse, error)
}

type carServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCarServiceClient(cc grpc.ClientConnInterface) CarServiceClient {
	return &carServiceClient{cc}
}

func (c *carServiceClient) CreateCar(ctx context.Context, in *CreateCarRequest, opts ...grpc.CallOption) (*CarEntity, error) {
	out := new(CarEntity)
	err := c.cc.Invoke(ctx, "/car.v1.CarService/CreateCar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *carServiceClient) GetCar(ctx context.Context, in *GetCarRequest, opts ...grpc.CallOption) (*Car, error) {
	out := new(Car)
	err := c.cc.Invoke(ctx, "/car.v1.CarService/GetCar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *carServiceClient) GetCars(ctx context.Context, in *GetCarsRequest, opts ...grpc.CallOption) (*GetCarsResponse, error) {
	out := new(GetCarsResponse)
	err := c.cc.Invoke(ctx, "/car.v1.CarService/GetCars", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *carServiceClient) LockCar(ctx context.Context, in *LockCarRequest, opts ...grpc.CallOption) (*LockCarResponse, error) {
	out := new(LockCarResponse)
	err := c.cc.Invoke(ctx, "/car.v1.CarService/LockCar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *carServiceClient) UnLockCar(ctx context.Context, in *UnLockCarRequest, opts ...grpc.CallOption) (*UnLockCarResponse, error) {
	out := new(UnLockCarResponse)
	err := c.cc.Invoke(ctx, "/car.v1.CarService/UnLockCar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *carServiceClient) UpdateCar(ctx context.Context, in *UpdateCarRequest, opts ...grpc.CallOption) (*UpdateCarResponse, error) {
	out := new(UpdateCarResponse)
	err := c.cc.Invoke(ctx, "/car.v1.CarService/UpdateCar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CarServiceServer is the server API for CarService service.
// All implementations must embed UnimplementedCarServiceServer
// for forward compatibility
type CarServiceServer interface {
	CreateCar(context.Context, *CreateCarRequest) (*CarEntity, error)
	GetCar(context.Context, *GetCarRequest) (*Car, error)
	GetCars(context.Context, *GetCarsRequest) (*GetCarsResponse, error)
	LockCar(context.Context, *LockCarRequest) (*LockCarResponse, error)
	UnLockCar(context.Context, *UnLockCarRequest) (*UnLockCarResponse, error)
	UpdateCar(context.Context, *UpdateCarRequest) (*UpdateCarResponse, error)
	mustEmbedUnimplementedCarServiceServer()
}

// UnimplementedCarServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCarServiceServer struct {
}

func (UnimplementedCarServiceServer) CreateCar(context.Context, *CreateCarRequest) (*CarEntity, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCar not implemented")
}
func (UnimplementedCarServiceServer) GetCar(context.Context, *GetCarRequest) (*Car, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCar not implemented")
}
func (UnimplementedCarServiceServer) GetCars(context.Context, *GetCarsRequest) (*GetCarsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCars not implemented")
}
func (UnimplementedCarServiceServer) LockCar(context.Context, *LockCarRequest) (*LockCarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LockCar not implemented")
}
func (UnimplementedCarServiceServer) UnLockCar(context.Context, *UnLockCarRequest) (*UnLockCarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnLockCar not implemented")
}
func (UnimplementedCarServiceServer) UpdateCar(context.Context, *UpdateCarRequest) (*UpdateCarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCar not implemented")
}
func (UnimplementedCarServiceServer) mustEmbedUnimplementedCarServiceServer() {}

// UnsafeCarServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CarServiceServer will
// result in compilation errors.
type UnsafeCarServiceServer interface {
	mustEmbedUnimplementedCarServiceServer()
}

func RegisterCarServiceServer(s grpc.ServiceRegistrar, srv CarServiceServer) {
	s.RegisterService(&CarService_ServiceDesc, srv)
}

func _CarService_CreateCar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CarServiceServer).CreateCar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/car.v1.CarService/CreateCar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CarServiceServer).CreateCar(ctx, req.(*CreateCarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CarService_GetCar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CarServiceServer).GetCar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/car.v1.CarService/GetCar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CarServiceServer).GetCar(ctx, req.(*GetCarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CarService_GetCars_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCarsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CarServiceServer).GetCars(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/car.v1.CarService/GetCars",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CarServiceServer).GetCars(ctx, req.(*GetCarsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CarService_LockCar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LockCarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CarServiceServer).LockCar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/car.v1.CarService/LockCar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CarServiceServer).LockCar(ctx, req.(*LockCarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CarService_UnLockCar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnLockCarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CarServiceServer).UnLockCar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/car.v1.CarService/UnLockCar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CarServiceServer).UnLockCar(ctx, req.(*UnLockCarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CarService_UpdateCar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CarServiceServer).UpdateCar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/car.v1.CarService/UpdateCar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CarServiceServer).UpdateCar(ctx, req.(*UpdateCarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CarService_ServiceDesc is the grpc.ServiceDesc for CarService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CarService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "car.v1.CarService",
	HandlerType: (*CarServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCar",
			Handler:    _CarService_CreateCar_Handler,
		},
		{
			MethodName: "GetCar",
			Handler:    _CarService_GetCar_Handler,
		},
		{
			MethodName: "GetCars",
			Handler:    _CarService_GetCars_Handler,
		},
		{
			MethodName: "LockCar",
			Handler:    _CarService_LockCar_Handler,
		},
		{
			MethodName: "UnLockCar",
			Handler:    _CarService_UnLockCar_Handler,
		},
		{
			MethodName: "UpdateCar",
			Handler:    _CarService_UpdateCar_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "car.proto",
}
