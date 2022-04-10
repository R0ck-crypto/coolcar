package dao

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	mongotesting "coolcar/shared/mongo/testing"
	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/testing/protocmp"
	"os"
	"testing"
)

func TestMongo_CreateTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb:%v", err)
	}

	db := mc.Database("coolcar")
	err = mongotesting.SetupIndexes(c, db)
	if err != nil {
		t.Fatalf("cannot setup indexes:%v", err)
	}

	m := NewMongo(db)

	cases := []struct {
		name       string
		tripID     string
		accountID  string
		tripStatus rentalpb.TripStatus
		wantErr    bool
	}{
		{
			name:       "finished",
			tripID:     "621e13d9b7e28adbf1aadbe0",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_FINISHED,
		},
		{
			name:       "another_finished",
			tripID:     "621e13d9b7e28adbf1aadbe1",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_FINISHED,
		},
		{
			name:       "in_progress",
			tripID:     "621e13d9b7e28adbf1aadbe2",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
		},
		{
			name:       "in_progress",
			tripID:     "621e13d9b7e28adbf1aadbe3",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
			wantErr:    true,
		},
		{
			name:       "in_progress_by_another_account",
			tripID:     "621e13d9b7e28adbf1aadbe4",
			accountID:  "account2",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
		},
	}
	for _, cc := range cases {
		mgutil.NewObjIDWithValue(id.TripID(cc.tripID))
		tr, err := m.CreateTrip(c, &rentalpb.Trip{
			AccountId: cc.accountID,
			Status:    cc.tripStatus,
		})
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s:error expected;got none", cc.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s:error creating trip:%v", cc.name, err)
			continue
		}
		if tr.ID.Hex() != cc.tripID {
			t.Errorf("%s:incorrect trip id;want:%q; got:%q",
				cc.name, cc.tripID, tr.ID.Hex())
		}

	}
}
func TestMongo_GetTrip(t *testing.T) {

	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb:%v", err)
	}

	m := NewMongo(mc.Database("coolcar"))
	acct := id.AccountID("account1")
	mgutil.NewObjID = primitive.NewObjectID
	tripRecord, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: acct.String(),
		CarID:     "car1",
		Start: &rentalpb.LocationStaus{
			PoiName: "startpoint",
			Location: &rentalpb.Location{
				Latitude:  30,
				Longitude: 120,
			},
		},
		End: &rentalpb.LocationStaus{
			PoiName:  "endpoint",
			FeeCent:  10000,
			KmDriven: 35,
			Location: &rentalpb.Location{
				Latitude:  35,
				Longitude: 115,
			},
		},
		Status: rentalpb.TripStatus_FINISHED,
	})
	if err != nil {
		t.Errorf("cannot create trip:%v", err)
	}

	got, err := m.GetTrip(c, objid.ToTripID(tripRecord.ID), acct)
	if err != nil {
		t.Errorf("got trip:%+v", got)
	}

	opt := protocmp.Transform()
	if diff := cmp.Diff(tripRecord, got, opt); diff != "" {
		t.Errorf("results differs; -want +got:%s", diff)
	}

}

func TestMongo_GetTrips(t *testing.T) {
	rows := []struct {
		id        string
		accountID string
		status    rentalpb.TripStatus
	}{
		{

			id:        "631e13d9b7e28adbf1aadbe0",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{

			id:        "631e13d9b7e28adbf1aadbe1",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{

			id:        "631e13d9b7e28adbf1aadbe2",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_IN_PROGRESS,
		},
		{

			id:        "631e13d9b7e28adbf1aadbe3",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{

			id:        "621e13d9b7e28adbf1aadbe4",
			accountID: "account_id_for_get_trips_1",
			status:    rentalpb.TripStatus_IN_PROGRESS,
		},
	}

	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	m := NewMongo(mc.Database("coolcar"))

	if err != nil {
		t.Fatalf("cannot connect mongodb:%v", err)
	}

	for _, r := range rows {
		mgutil.NewObjIDWithValue(id.TripID(r.id))
		m.CreateTrip(c, &rentalpb.Trip{
			AccountId: r.accountID,
			Status:    r.status,
		})
		if err != nil {
			t.Fatalf("cannot create rows:%v", err)
		}
	}

	cases := []struct {
		name       string
		accountID  string
		status     rentalpb.TripStatus
		wantCount  int
		wantOnlyID string
	}{
		{
			name:      "get_all",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_TS_NOT_SPECIFIED,
			wantCount: 4,
		},
		{
			name:       "get_in_progress",
			accountID:  "account_id_for_get_trips",
			status:     rentalpb.TripStatus_IN_PROGRESS,
			wantCount:  1,
			wantOnlyID: "631e13d9b7e28adbf1aadbe2",
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			records, err := m.GetTrips(context.Background(),
				id.AccountID(cc.accountID),
				cc.status)
			if err != nil {
				t.Errorf("cannot get trips:%v", err)
			}

			if cc.wantCount != len(records) {
				t.Errorf("incorrect result count;want:%d, got:%d",
					cc.wantCount, len(records))
			}

			if cc.wantOnlyID != "" && len(records) > 0 {
				if cc.wantOnlyID != records[0].ID.Hex() {
					t.Errorf("only_id incorrect; want:%q,got:%q",
						cc.wantOnlyID, records[0].ID.Hex())
				}
			}
		})
	}
}

func TestMongo_UpdateTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb:%v", err)
	}

	m := NewMongo(mc.Database("coolcar"))
	tid := id.TripID("631f13d9b7e28adbf1aadbe2")
	aid := id.AccountID("account_for_update")

	var now int64 = 10000
	mgutil.NewObjIDWithValue(tid)
	mgutil.UpdateAt = func() int64 {
		return now
	}
	tripRecord, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: aid.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStaus{
			PoiName: "start_poi",
		},
	})
	if err != nil {
		t.Fatalf("cannot create trip:%v", err)
	}

	if tripRecord.UpdateAt != 10000 {
		t.Fatalf("wrong updated at; want:10000;got:%d", tripRecord.UpdateAt)
	}

	updates := &rentalpb.Trip{
		AccountId: aid.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStaus{
			PoiName: "start_poi_updated",
		},
	}

	cases := []struct {
		name          string
		now           int64
		withUpdatedAt int64
		wantErr       bool
	}{
		{
			name:          "normal_update",
			now:           20000,
			withUpdatedAt: 10000,
		},
		{
			name:          "update_with_scale_timestamp",
			now:           30000,
			withUpdatedAt: 10000,
			wantErr:       true,
		},
		{
			name:          "update_with_refetch",
			now:           40000,
			withUpdatedAt: 20000,
		},
	}

	for _, cc := range cases {
		now = cc.now
		err := m.UpdateTrip(c, tid, aid, cc.withUpdatedAt, updates)
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s:want error; got none", cc.name)
			} else {
				continue
			}
		} else {
			if err != nil {
				t.Errorf("%s:cannot update:%v", cc.name, err)
			}
		}

		updatedtrip, err := m.GetTrip(c, tid, aid)
		if err != nil {
			t.Errorf("%s:cannot get trip after update:%v", cc.name, err)
		}

		if now != updatedtrip.UpdateAt {
			t.Errorf("%s:incorrect updatedAt: want:%d, got:%d",
				cc.name, now, updatedtrip.UpdateAt)
		}

	}

}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
