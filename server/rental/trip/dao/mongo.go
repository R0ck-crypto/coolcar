package dao

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	tripField      = "trip"
	accountIDField = tripField + ".accountid"
	statusField    = tripField + ".status"
)

// Mongo define a mongo dao
type Mongo struct {
	col *mongo.Collection
}

// NewMongo creates a mongo dao
func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col: db.Collection("trip"),
	}
}

type TripRecord struct {
	mgutil.IDFiled       `bson:"inline"`
	mgutil.UpdateAtField `bson:"inline"`
	Trip                 *rentalpb.Trip `bson:"trip"`
}

// CreateTrip creates a trip
func (m *Mongo) CreateTrip(c context.Context, trip *rentalpb.Trip) (*TripRecord, error) {
	r := &TripRecord{
		Trip: trip,
	}
	r.ID = mgutil.NewObjID()
	r.UpdateAt = mgutil.UpdateAt()

	_, err := m.col.InsertOne(c, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (m *Mongo) GetTrip(c context.Context, id id.TripID, accountID id.AccountID) (*TripRecord, error) {
	objID, err := objid.FromID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id:%v", err)
	}
	result := m.col.FindOne(c, bson.M{
		mgutil.IDField: objID,
		accountIDField: accountID,
	})
	if err := result.Err(); err != nil {
		return nil, err
	}

	var tr TripRecord
	err = result.Decode(&tr)
	if err != nil {
		return nil, fmt.Errorf("can not decode:%v", err)
	}

	return &tr, nil

}

func (m *Mongo) GetTrips(c context.Context, accountID id.AccountID, stautus rentalpb.TripStatus) ([]*TripRecord, error) {
	filter := bson.M{
		accountIDField: accountID.String(),
	}
	if stautus != rentalpb.TripStatus_TS_NOT_SPECIFIED {
		filter[statusField] = stautus
	}

	res, err := m.col.Find(c, filter)
	if err != nil {
		return nil, err
	}

	var trips []*TripRecord
	for res.Next(c) {
		var trip TripRecord
		err := res.Decode(&trip)
		if err != nil {
			return nil, err
		}
		trips = append(trips, &trip)
	}
	return trips, nil

}

func (m *Mongo) UpdateTrip(c context.Context, tripID id.TripID, accountID id.AccountID, updatedAt int64, trip *rentalpb.Trip) error {
	objectID, err := objid.FromID(tripID)
	if err != nil {
		return fmt.Errorf("invalid id :%v", err)
	}

	newUpdateAt := mgutil.UpdateAt()
	res, err := m.col.UpdateOne(c, bson.M{
		mgutil.IDField:            objectID,
		accountIDField:            accountID,
		mgutil.UpdatedAtFieldName: updatedAt,
	}, mgutil.Set(bson.M{
		tripField:                 trip,
		mgutil.UpdatedAtFieldName: newUpdateAt,
	}))

	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil

}
