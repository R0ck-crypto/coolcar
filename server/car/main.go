package main

import (
	"context"
	"coolcar/car/amqpclt"
	carpb "coolcar/car/api/gen/v1"
	"coolcar/car/car"
	"coolcar/car/dao"
	"coolcar/shared/server"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
)

func main() {

	logger, err := server.NewZapLogger()

	if err != nil {
		log.Fatalf("cannot create logger:%v", logger)
	}

	c := context.Background()
	mc, err := mongo.Connect(c,
		options.Client().ApplyURI("mongodb://admin:Hacker007@192.168.2.196:27017"))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}

	db := mc.Database("coolcar")

	conn, err := amqp.Dial("amqp://guest:guest@192.168.2.196:5672")
	if err != nil {
		logger.Fatal("cannot dial amqp")
	}

	pub, err := amqpclt.NewPublisher(conn, "coolcar")
	if err != nil {
		logger.Fatal("cannot create publisher")
	}

	logger.Sugar().Fatal(server.RunGRPCServer(
		&server.GRPConfig{
			Name:   "car",
			Addr:   ":8084",
			Logger: logger,
			RegisterFunc: func(s *grpc.Server) {
				carpb.RegisterCarServiceServer(s, &car.Service{
					Logger:    logger,
					Mongo:     dao.NewMongo(db),
					Publisher: pub,
				})
			},
		},
	))

}
