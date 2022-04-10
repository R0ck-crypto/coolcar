package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/dao"
	"coolcar/auth/token"
	"coolcar/auth/wechat"
	"coolcar/shared/server"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
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

	pkFile, err := os.Open("auth/private.key")
	if err != nil {
		log.Fatal("cannot open private key", zap.Error(err))
	}

	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		log.Fatal(" cannot read private key", zap.Error(err))
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		log.Fatal(" cannot parse private key", zap.Error(err))
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPConfig{
		Name:   "auth",
		Addr:   ":8081",
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			authpb.RegisterAuthServiceServer(s, &auth.Service{
				OpenIDResolver: &wechat.Service{
					AppID:     "wxf3530dbf1a8ddc79",
					AppSecret: "010f43c19045f2df230a482b36093def",
				},
				Mongo:          dao.NewMongo(mc.Database("coolcar")),
				Logger:         logger,
				TokenExpire:    2 * time.Hour,
				TokenGenerator: token.NewJwtTokeGen("coolcar/auth", privateKey),
			})
		},
	}))

}
