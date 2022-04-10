package poi

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"google.golang.org/protobuf/proto"
	"hash/fnv"
)

var poi = []string{
	"中关村",
	"天安门",
	"陆家嘴",
	"迪士尼",
	"天河体育中心",
	"广州塔",
}

type Manager struct {
}

func (m *Manager) Resolve(ctx context.Context, location *rentalpb.Location) (string, error) {
	b, err := proto.Marshal(location)
	if err != nil {
		return "", nil
	}
	h := fnv.New32()
	h.Write(b)
	return poi[int(h.Sum32())%len(poi)], nil
}
