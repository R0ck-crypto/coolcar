package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/client/poi"
	"coolcar/rental/trip/dao"
	token "coolcar/shared/auth"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	mongotesting "coolcar/shared/mongo/testing"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"testing"
)

func TestService_CreateTrip(t *testing.T) {
	c := context.Background()

	pm := &profileManager{}
	cm := &carManager{}
	s := NewService(c, t, pm, cm)

	req := &rentalpb.CreateTripRequest{
		CarId: "car1",
		Start: &rentalpb.Location{
			Latitude:  32.1,
			Longitude: 114.1,
		},
	}

	pm.iID = "identity1"
	golden := `{"account_id":%q,"carID":"car1","start":{"location":{"latitude":32.1,"longitude":114.1},"poi_name":"天河体育中心","timestamp_sec":1647680506},"current":{"location":{"latitude":32.1,"longitude":114.1},"poi_name":"天河体育中心","timestamp_sec":1647680506},"status":1,"identity_id":"identity1"}`
	nowFunc = func() int64 {
		return 1647680506
	}
	cases := []struct {
		name         string
		accountId    string
		tripid       string
		profileErr   error
		carVerifyErr error
		carUnlockErr error
		want         string
		wantErr      bool
	}{
		{
			name:      "normal_create",
			accountId: "account1",
			tripid:    "631f13d9b7e28adbf1aadbe4",
			want:      fmt.Sprintf(golden, "account1"),
		},
		{
			name:       "profile_err",
			accountId:  "account2",
			tripid:     "631f13d9b7e28adbf1aadbe5",
			profileErr: fmt.Errorf("profile"),
			wantErr:    true,
		},
		{
			name:         "car_verify_err",
			accountId:    "account3",
			tripid:       "631f13d9b7e28adbf1aadbe6",
			carVerifyErr: fmt.Errorf("verify"),
			wantErr:      true,
		},
		{
			name:         "car_unlock_err",
			accountId:    "account4",
			tripid:       "631f13d9b7e28adbf1aadbe7",
			carUnlockErr: fmt.Errorf("unlock"),
			want:         fmt.Sprintf(golden, "account4"),
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			mgutil.NewObjIDWithValue(id.TripID(cc.tripid))
			pm.err = cc.profileErr
			cm.verifyErr = cc.carVerifyErr
			cm.unlockErr = cc.carUnlockErr
			c = token.ContextWithAccountID(context.Background(), id.AccountID(cc.accountId))
			res, err := s.CreateTrip(c, req)
			if cc.wantErr {
				if err == nil {
					t.Errorf("want error; got none")
				} else {
					return
				}
			}

			if err != nil {
				t.Errorf("error creating trip:%v", err)
				return
			}

			if res.Id != cc.tripid {
				t.Errorf("incorrect id,want:%q, got:%q", cc.tripid, res.Id)
			}

			bytes, err := json.Marshal(res.Trip)
			if err != nil {
				t.Errorf("can not marshall response:%v", err)
			}

			tripstr := string(bytes)

			if cc.want != tripstr {
				t.Errorf("incorrect response,want:%s, got:%s", cc.want, tripstr)
			}
		})
	}
}

func NewService(c context.Context, t *testing.T, pm ProfileManager, cm CarManager) *Service {

	client, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf(" cannot create mongo client:%v", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("cannot create logger:%v", err)
	}

	db := client.Database("coolcar")
	mongotesting.SetupIndexes(c, db)
	return &Service{
		ProfileManager: pm,
		CarManager:     cm,
		POIManager:     &poi.Manager{},
		Mongo:          dao.NewMongo(db),
		Logger:         logger,
	}
}

func TestLifeCycle(t *testing.T) {

	c := token.ContextWithAccountID(context.Background(), id.AccountID("id_for_life_cycle"))

	s := NewService(c, t, &profileManager{}, &carManager{})

	tid := id.TripID("633f13d9b7e28adbf1aadbe4")
	mgutil.NewObjIDWithValue(tid)

	cases := []struct {
		name string
		now  int64
		op   func() (*rentalpb.Trip, error)
		want string
	}{
		{
			name: "create_trip",
			now:  10000,
			op: func() (*rentalpb.Trip, error) {
				e, err := s.CreateTrip(c, &rentalpb.CreateTripRequest{
					CarId: "car1",
					Start: &rentalpb.Location{
						Latitude:  32.1,
						Longitude: 114.1,
					},
				})
				if err != nil {
					return nil, err
				}

				return e.Trip, nil
			},
			want: `{"account_id":"id_for_life_cycle","carID":"car1","start":{"location":{"latitude":32.1,"longitude":114.1},"poi_name":"天河体育中心","timestamp_sec":10000},"current":{"location":{"latitude":32.1,"longitude":114.1},"poi_name":"天河体育中心","timestamp_sec":10000},"status":1}`,
		},
		{
			name: "update_trip",
			now:  20000,
			op: func() (*rentalpb.Trip, error) {
				return s.UpdateTrip(c, &rentalpb.UpdateTripRequest{
					Id: tid.String(),
					Current: &rentalpb.Location{
						Latitude:  28.2345,
						Longitude: 123.15454,
					},
				})
			},
			want: `{"account_id":"id_for_life_cycle","carID":"car1","start":{"location":{"latitude":32.1,"longitude":114.1},"poi_name":"天河体育中心","timestamp_sec":10000},"current":{"location":{"latitude":28.2345,"longitude":123.15454},"fee_cent":9116,"km_driven":359.56460921309167,"poi_name":"广州塔","timestamp_sec":20000},"status":1}`,
		},
		{
			name: "finish_trip",
			now:  30000,
			op: func() (*rentalpb.Trip, error) {
				return s.UpdateTrip(c, &rentalpb.UpdateTripRequest{
					Id:      tid.String(),
					EndTrip: true,
				})
			},
			want: `{"account_id":"id_for_life_cycle","carID":"car1","start":{"location":{"latitude":32.1,"longitude":114.1},"poi_name":"天河体育中心","timestamp_sec":10000},"current":{"location":{"latitude":28.2345,"longitude":123.15454},"fee_cent":17638,"km_driven":480.5067823692563,"poi_name":"广州塔","timestamp_sec":30000},"end":{"location":{"latitude":28.2345,"longitude":123.15454},"fee_cent":17638,"km_driven":480.5067823692563,"poi_name":"广州塔","timestamp_sec":30000},"status":2}`,
		},
		{
			name: "query_trip",
			now:  40000,
			op: func() (*rentalpb.Trip, error) {
				return s.GetTrip(c, &rentalpb.GetTripRequest{
					Id: tid.String(),
				})
			},
			want: `{"account_id":"id_for_life_cycle","carID":"car1","start":{"location":{"latitude":32.1,"longitude":114.1},"poi_name":"天河体育中心","timestamp_sec":10000},"current":{"location":{"latitude":28.2345,"longitude":123.15454},"fee_cent":17638,"km_driven":480.5067823692563,"poi_name":"广州塔","timestamp_sec":30000},"end":{"location":{"latitude":28.2345,"longitude":123.15454},"fee_cent":17638,"km_driven":480.5067823692563,"poi_name":"广州塔","timestamp_sec":30000},"status":2}`,
		},
	}

	rand.Seed(1234)
	for _, cc := range cases {
		nowFunc = func() int64 {
			return cc.now
		}
		trip, err := cc.op()
		if err != nil {
			t.Errorf("%s:operation failed:%v", cc.name, err)
			continue
		}

		bytes, err := json.Marshal(trip)
		if err != nil {
			t.Errorf("%s:failed marshalling response:%v", cc.name, err)
		}

		got := string(bytes)

		if cc.want != got {
			t.Errorf("%s: incorrect response; want:%s, got:%s", cc.name, cc.want, got)
		}

	}

}

type profileManager struct {
	iID id.IdentityID
	err error
}

func (p *profileManager) Verify(ctx context.Context, accountID id.AccountID) (id.IdentityID, error) {
	return p.iID, p.err
}

type carManager struct {
	verifyErr error
	unlockErr error
}

func (c carManager) Verify(ctx context.Context, carID id.CarID, location *rentalpb.Location) error {
	return c.verifyErr
}

func (c carManager) Unlock(ctx context.Context, carID id.CarID, aid id.AccountID, tid id.TripID, avatarURL string) error {
	return c.unlockErr
}

func (c carManager) Lock(ctx context.Context, carID id.CarID) error {
	return nil
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
