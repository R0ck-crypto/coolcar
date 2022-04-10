package car

import (
	"context"
	carpb "coolcar/car/api/gen/v1"
	"coolcar/car/dao"
	"coolcar/shared/id"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Publisher interface {
	Publish(ctx context.Context, entity *carpb.CarEntity) error
}

// Service defines a car service.
type Service struct {
	Mongo     *dao.Mongo
	Logger    *zap.Logger
	Publisher Publisher
	carpb.UnimplementedCarServiceServer
}

// CreateCar creates a car.
func (s Service) CreateCar(ctx context.Context, request *carpb.CreateCarRequest) (*carpb.CarEntity, error) {
	cr, err := s.Mongo.CreateCar(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &carpb.CarEntity{
		Id:  cr.ID.Hex(),
		Car: cr.Car,
	}, nil
}

// GetCar gets car.
func (s Service) GetCar(ctx context.Context, request *carpb.GetCarRequest) (*carpb.Car, error) {
	cr, err := s.Mongo.GetCar(ctx, id.CarID(request.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "")
	}
	return cr.Car, nil
}

// GetCars gets cars.
func (s Service) GetCars(ctx context.Context, request *carpb.GetCarsRequest) (*carpb.GetCarsResponse, error) {
	cars, err := s.Mongo.GetCars(ctx)
	if err != nil {
		s.Logger.Error("cannot get cars", zap.Error(err))
		return nil, err
	}
	res := &carpb.GetCarsResponse{}

	for _, cr := range cars {
		res.Cars = append(res.Cars, &carpb.CarEntity{
			Id:  cr.ID.Hex(),
			Car: cr.Car,
		})
	}
	return res, nil

}

func (s Service) LockCar(ctx context.Context, request *carpb.LockCarRequest) (*carpb.LockCarResponse, error) {
	car, err := s.Mongo.UpdateCar(ctx, id.CarID(request.Id), carpb.CarStatus_UNLOCKED, &dao.CarUpdate{
		Status: carpb.CarStatus_LOCKING,
	})
	if err != nil {
		code := codes.Internal
		if err == mongo.ErrNoDocuments {
			code = codes.NotFound
		}
		return nil, status.Errorf(code, "cannot updpate:%v", err)
	}

	s.publish(ctx, car)
	return &carpb.LockCarResponse{}, nil
}

//
func (s Service) UnLockCar(ctx context.Context, request *carpb.UnLockCarRequest) (*carpb.UnLockCarResponse, error) {
	car, err := s.Mongo.UpdateCar(ctx, id.CarID(request.Id), carpb.CarStatus_LOCKED, &dao.CarUpdate{
		Status:       carpb.CarStatus_UNLOCKING,
		Driver:       request.Driver,
		UpdateTripID: true,
		TripId:       id.TripID(request.TripId),
	})
	if err != nil {
		code := codes.Internal
		if err == mongo.ErrNoDocuments {
			code = codes.NotFound
		}
		return nil, status.Errorf(code, "cannot update:%v", err)
	}
	s.publish(ctx, car)
	return &carpb.UnLockCarResponse{}, nil
}

// UpdateCar updates the cars position and status without
func (s Service) UpdateCar(ctx context.Context, request *carpb.UpdateCarRequest) (*carpb.UpdateCarResponse, error) {

	update := &dao.CarUpdate{
		Status:   request.Status,
		Position: request.Position,
	}

	if request.Status == carpb.CarStatus_LOCKED {
		update.Driver = &carpb.Driver{}
		update.UpdateTripID = true
		update.TripId = id.TripID("")
	}

	car, err := s.Mongo.UpdateCar(ctx, id.CarID(request.Id), carpb.CarStatus_CS_NOT_SPECIFIED, update)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.publish(ctx, car)

	return &carpb.UpdateCarResponse{}, nil

}

func (s Service) publish(ctx context.Context, car *dao.CarRecord) {
	err := s.Publisher.Publish(ctx, &carpb.CarEntity{
		Id:  car.ID.Hex(),
		Car: car.Car,
	})
	if err != nil {
		s.Logger.Warn("cannot publish", zap.Error(err))
	}

}
