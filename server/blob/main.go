package main

import (
	"context"
	blobbpb "coolcar/blob/api/gen/v1"
	"coolcar/blob/blob"
	"coolcar/blob/cos"
	"coolcar/blob/dao"
	"coolcar/shared/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
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
	st, err := cos.NewService("https://coolcar-1310568949.cos.ap-guangzhou.myqcloud.com",
		"AKIDwlLPJyNXT9ZjeIfla9RSCguTnsgCSQ8a",
		"C7mAG5XmkRr2wZv0pNfkY5LMAhyzkERO",
	)
	if err != nil {
		logger.Fatal("cannot create cos service")
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPConfig{
		Name:   "blob",
		Addr:   ":8083",
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			blobbpb.RegisterBlobServiceServer(s, &blob.Service{
				Mongo:   dao.NewMongo(db),
				Logger:  logger,
				Storage: st,
			})
		},
	}))

}
