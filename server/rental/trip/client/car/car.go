package car

import (
	"context"
	carpb "coolcar/car/api/gen/v1"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
	"fmt"
)

// Manager defines a car manager
type Manager struct {
	CarService carpb.CarServiceClient
}

// Verify verifies car status
func (m Manager) Verify(ctx context.Context, carID id.CarID, location *rentalpb.Location) error {
	car, err := m.CarService.GetCar(ctx, &carpb.GetCarRequest{
		Id: carID.String(),
	})
	if err != nil {
		return fmt.Errorf("cannot get car :v", err)
	}

	if car.Status != carpb.CarStatus_LOCKED {
		return fmt.Errorf("cannot unlock; car status is %v", car.Status)
	}

	return nil
}

// Unlock unlocks a car
func (m Manager) Unlock(ctx context.Context, carID id.CarID, aid id.AccountID, tid id.TripID, avatarURL string) error {
	_, err := m.CarService.UnLockCar(ctx, &carpb.UnLockCarRequest{
		Id: carID.String(),
		Driver: &carpb.Driver{
			Id:        aid.String(),
			AvatarUrl: avatarURL,
		},
		TripId: tid.String(),
	})
	if err != nil {
		return fmt.Errorf("cannot unlock:%v", err)
	}
	return nil
}

func (m Manager) Lock(ctx context.Context, carID id.CarID) error {
	_, err := m.CarService.LockCar(ctx, &carpb.LockCarRequest{
		Id: carID.String(),
	})
	if err != nil {
		return fmt.Errorf("cannot lock:%v", err)
	}
	return nil
}
