package profile

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
	"encoding/base64"
	"fmt"
	"google.golang.org/protobuf/proto"
)

type Fetcher interface {
	GetProfile(ctx context.Context, request *rentalpb.GetProfileRequest) (*rentalpb.Profile, error)
}

type Manager struct {
	Fetcher Fetcher
}

func (m Manager) Verify(ctx context.Context, accountID id.AccountID) (id.IdentityID, error) {
	nilId := id.IdentityID("")

	profile, err := m.Fetcher.GetProfile(ctx, &rentalpb.GetProfileRequest{})
	if err != nil {
		return nilId, fmt.Errorf("can not get profile:%v", err)
	}

	if profile.IdentityStatus != rentalpb.IdentityStatus_VERIFIED {
		return nilId, fmt.Errorf("invalid indentity status")
	}

	bytes, err := proto.Marshal(profile.Identity)
	if err != nil {
		return nilId, fmt.Errorf("can not marshal identity: %v", err)
	}

	return id.IdentityID(base64.StdEncoding.EncodeToString(bytes)), nil

}
