package main

import (
	"context"
	blobpb "coolcar/blob/api/gen/v1"
	carpb "coolcar/car/api/gen/v1"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/profile"
	profiledao "coolcar/rental/profile/dao"
	"coolcar/rental/trip"
	"coolcar/rental/trip/client/car"
	"coolcar/rental/trip/client/poi"
	profCient "coolcar/rental/trip/client/profile"
	tripdao "coolcar/rental/trip/dao"
	"coolcar/shared/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot use logger:%v", err)
	}

	c := context.Background()
	mc, err := mongo.Connect(c,
		options.Client().ApplyURI("mongodb://admin:Hacker007@192.168.2.196:27017"))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}

	db := mc.Database("coolcar")

	blobConn, err := grpc.Dial("localhost:8083", grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect blob service", zap.Error(err))
	}

	profService := &profile.Service{
		BlobClient:        blobpb.NewBlobServiceClient(blobConn),
		PhotoGetExpire:    4 * time.Second,
		PhotoUploadExpire: 10 * time.Second,
		Mongo:             profiledao.NewMongo(db),
		Logger:            logger,
	}

	carConn, err := grpc.Dial("localhost:8084", grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect car service", zap.Error(err))
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPConfig{
		Name:              "rental",
		Addr:              ":8082",
		Logger:            logger,
		AuthPublickeyfile: "shared/auth/public.key",
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				CarManager: &car.Manager{
					CarService: carpb.NewCarServiceClient(carConn),
				},

				ProfileManager: &profCient.Manager{
					Fetcher: profService,
				},
				POIManager: &poi.Manager{},
				Mongo:      tripdao.NewMongo(db),
				Logger:     logger,
			})
			rentalpb.RegisterProfileServiceServer(s, profService)
		},
	}))

}
