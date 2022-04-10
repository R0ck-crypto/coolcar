package dao

import (
	"context"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	mongotesting "coolcar/shared/mongo/testing"
	"go.mongodb.org/mongo-driver/bson"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

var mongoURI string

func TestMongo_ResolveAccountID(t *testing.T) {

	// create container
	c := context.Background()
	mc, err := mongo.Connect(c,
		options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("cannot connect mongodb:%v", err)
	}
	//col := mc.Database("coolcar").Collection("account")
	m := NewMongo(mc.Database("coolcar"))
	_, err = m.col.InsertMany(c, []interface{}{
		bson.M{
			mgutil.IDField: objid.MustFromID(id.AccountID("61fc9bb34134d8d76da79d4e")),
			openIDField:    "openid_1",
		},
		bson.M{
			mgutil.IDField: objid.MustFromID(id.AccountID("61fc9bb34134d8d76da79d4f")),
			openIDField:    "openid_2",
		},
	})
	if err != nil {
		t.Fatalf("cannot insert initial values:%v", err)
	}

	mgutil.NewObjIDWithValue(id.AccountID("61fc9bb34134d8d76da79d4d"))

	cases := []struct {
		name    string
		open_id string
		want    string
	}{
		{
			name:    "existing_user",
			open_id: "openid_1",
			want:    "61fc9bb34134d8d76da79d4e",
		},
		{
			name:    "another_existing_user",
			open_id: "openid_2",
			want:    "61fc9bb34134d8d76da79d4f",
		},
		{
			name:    "new_user",
			open_id: "openid_3",
			want:    "61fc9bb34134d8d76da79d4d",
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			id, err := m.ResolveAccountID(context.Background(), cc.open_id)
			if err != nil {
				t.Errorf("fail resolve account id for 123:%v", err)
			}
			if id.String() != cc.want {
				t.Errorf("resolve account id: want: %q, got:%q", cc.want, id)
			}
		})
	}

	//id, err := m.ResolveAccountID(c, "123")
	//if err != nil {
	//	t.Errorf("fail resolve account id for 123:%v", err)
	//} else {
	//	want := "61fc9bb34134d8d76da79d4d"
	//	if id != want {
	//		t.Errorf("resolve account id: want: %q, got:%q", want, id)
	//	}
	//}
	//// remove container

}

func mustObjID(hex string) primitive.ObjectID {
	objID, _ := primitive.ObjectIDFromHex(hex)
	return objID
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
