package testing

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

const (
	image         = "mongo"
	containerPort = "27017/tcp"
)

var mongoURI string

const defaultMongoURI = "mongodb://admin:Hacker007@192.168.2.196:27017"

//RunWithMongoInDocker runs testing with
// a mongodb instance in a docker container
func RunWithMongoInDocker(m *testing.M) int {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation(),
		client.WithHost("tcp://192.168.2.196:2375"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
			ExposedPorts: nat.PortSet{
				containerPort: {},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				containerPort: []nat.PortBinding{
					{
						HostIP:   "192.168.2.196",
						HostPort: "0",
					},
				},
			},
		},
		nil,
		nil,
		"",
	)

	if err != nil {
		panic(err)
	}

	containerID := resp.ID
	defer func() {
		err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			panic(err)
		}

	}()

	err = cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	inspect, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		panic(err)
	}

	hostPort := inspect.NetworkSettings.Ports[containerPort][0]
	mongoURI = fmt.Sprintf("mongodb://%s:%s", hostPort.HostIP, hostPort.HostPort)

	return m.Run()
}

// New client create a client conneted to the mongo instantce
func NewClient(c context.Context) (*mongo.Client, error) {
	if mongoURI == "" {
		return nil, fmt.Errorf("mongo uri not set,Please run  RunWithMongoInDocker in TestMain")
	}
	return mongo.Connect(c, options.Client().ApplyURI(mongoURI))
}

func NewDefaultClient(c context.Context) (*mongo.Client, error) {
	if mongoURI == "" {
		return nil, fmt.Errorf("mongo uri not set,Please run  RunWithMongoInDocker in TestMain")
	}
	return mongo.Connect(c, options.Client().ApplyURI(defaultMongoURI))
}

func SetupIndexes(c context.Context, d *mongo.Database) error {
	_, err := d.Collection("account").Indexes().CreateOne(c, mongo.IndexModel{
		Keys: bson.D{
			{Key: "open_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		return err
	}

	_, err = d.Collection("trip").Indexes().CreateOne(c, mongo.IndexModel{
		Keys: bson.D{
			{Key: "trip.accountid", Value: 1},
			{Key: "trip.status", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
			"trip.status": 1,
		}),
	})

	_, err = d.Collection("profile").Indexes().CreateOne(c, mongo.IndexModel{
		Keys: bson.D{
			{Key: "accountid", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		return err
	}

	return err
}
