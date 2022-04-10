package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/dao"
	token "coolcar/shared/auth"
	"coolcar/shared/id"
	"coolcar/shared/mongo/objid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"time"
)

type Service struct {
	ProfileManager ProfileManager
	CarManager     CarManager
	POIManager     POIManager
	Mongo          *dao.Mongo
	Logger         *zap.Logger
	rentalpb.UnimplementedTripServiceServer
}

// ProfileManager  defines the ACL(Anti Corruption Layer)
// for profile verification logic
type ProfileManager interface {
	Verify(ctx context.Context, accountID id.AccountID) (id.IdentityID, error)
}

// CarManager define the ACL for car Management
type CarManager interface {
	Verify(ctx context.Context, carID id.CarID, location *rentalpb.Location) error
	Unlock(ctx context.Context, carID id.CarID, aid id.AccountID, tid id.TripID, avatarURL string) error
	Lock(ctx context.Context, carID id.CarID) error
}

// POIManager resolve POI(Point of interest)
type POIManager interface {
	Resolve(ctx context.Context, location *rentalpb.Location) (string, error)
}

func (s Service) CreateTrip(ctx context.Context, request *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {

	aid, err := token.AccountIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if request.CarId == "" || request.Start == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}

	// 验证驾驶者身份
	iID, err := s.ProfileManager.Verify(ctx, aid)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}

	// 检查车辆状态
	carID := id.CarID(request.CarId)
	err = s.CarManager.Verify(ctx, carID, request.Start)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}

	ls := s.calcCurrentStatus(ctx, &rentalpb.LocationStaus{
		Location:     request.Start,
		TimestampSec: nowFunc(),
	}, request.Start)

	trip, err := s.Mongo.CreateTrip(ctx, &rentalpb.Trip{
		AccountId:  aid.String(),
		CarID:      carID.String(),
		IdentityId: iID.String(),
		Status:     rentalpb.TripStatus_IN_PROGRESS,
		Start:      ls,
		Current:    ls,
	})
	if err != nil {
		s.Logger.Warn("cannot create trip", zap.Error(err))
		return nil, status.Error(codes.AlreadyExists, "")
	}

	// 车辆开锁
	go func() {
		err := s.CarManager.Unlock(context.Background(), carID, aid, objid.ToTripID(trip.ID), request.AvatarUrl)
		if err != nil {
			s.Logger.Error("cannot unlock car", zap.Error(err))
		}
	}()

	return &rentalpb.TripEntity{
		Id:   trip.ID.Hex(),
		Trip: trip.Trip,
	}, nil

}

func (s Service) GetTrip(ctx context.Context, request *rentalpb.GetTripRequest) (*rentalpb.Trip, error) {

	aid, err := token.AccountIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	tr, err := s.Mongo.GetTrip(ctx, id.TripID(request.Id), aid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "")
	}
	return tr.Trip, nil
}

func (s Service) GetTrips(ctx context.Context, request *rentalpb.GetTripsRequest) (*rentalpb.GetTripsReponse, error) {
	aid, err := token.AccountIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	trips, err := s.Mongo.GetTrips(ctx, aid, request.Status)
	if err != nil {
		s.Logger.Error("cannot get trips", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	res := &rentalpb.GetTripsReponse{}
	for _, tr := range trips {
		res.Trips = append(res.Trips, &rentalpb.TripEntity{
			Id:   tr.ID.Hex(),
			Trip: tr.Trip,
		})
	}
	return res, nil
}

// UpdateTrip updates a trip
func (s Service) UpdateTrip(ctx context.Context, request *rentalpb.UpdateTripRequest) (*rentalpb.Trip, error) {
	accountID, err := token.AccountIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "")
	}

	tid := id.TripID(request.Id)
	aid := id.AccountID(accountID)
	tr, err := s.Mongo.GetTrip(ctx, tid, aid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "")
	}

	if tr.Trip.Current == nil {
		s.Logger.Error("trip without current set", zap.String("id", tid.String()))
		return nil, status.Error(codes.Internal, "")
	}

	cur := tr.Trip.Current.Location
	if request.Current != nil {
		cur = request.Current
	}
	tr.Trip.Current = s.calcCurrentStatus(ctx, tr.Trip.Current, cur)

	if request.EndTrip {
		tr.Trip.End = tr.Trip.Current
		tr.Trip.Status = rentalpb.TripStatus_FINISHED
		err := s.CarManager.Lock(ctx, id.CarID(tr.Trip.CarID))
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "cannot lock car:%v", err)
		}
	}

	err = s.Mongo.UpdateTrip(ctx, tid, aid, tr.UpdateAt, tr.Trip)
	if err != nil {
		return nil, status.Error(codes.Aborted, "")
	}

	return tr.Trip, nil
}

var nowFunc = func() int64 {
	return time.Now().Unix()
}

const (
	centsPerSec = 0.7
	kmPerSec    = 0.02
)

func (s Service) calcCurrentStatus(ctx context.Context, last *rentalpb.LocationStaus, cur *rentalpb.Location) *rentalpb.LocationStaus {
	now := nowFunc()
	elapsedSec := float64(now - last.TimestampSec)
	poi, err := s.POIManager.Resolve(ctx, cur)
	if err != nil {
		s.Logger.Info("cannot resolve poi", zap.Stringer("location", cur), zap.Error(err))
	}
	return &rentalpb.LocationStaus{
		Location:     cur,
		FeeCent:      last.FeeCent + int32(centsPerSec*elapsedSec*2*rand.Float64()),
		KmDriven:     last.KmDriven + kmPerSec*elapsedSec*2*rand.Float64(),
		TimestampSec: now,
		PoiName:      poi,
	}

}

//func (s Service) mustEmbedUnimplementedTripServiceServer() {
//	panic("implement me")
//}
