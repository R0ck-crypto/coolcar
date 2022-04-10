package auth

import (
	"context"
	"coolcar/auth/dao"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	authpb "coolcar/auth/api/gen/v1"

	"go.uber.org/zap"
)

type Service struct {
	OpenIDResolver
	TokenGenerator
	TokenExpire time.Duration
	Logger      *zap.Logger
	Mongo       *dao.Mongo
	authpb.UnimplementedAuthServiceServer
}

type OpenIDResolver interface {
	Resolve(code string) (string, error)
}

type TokenGenerator interface {
	GenerateToken(accountID string, expire time.Duration) (string, error)
}

func (s *Service) Login(ctx context.Context, request *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	s.Logger.Info("received code", zap.String("code", request.Code))
	openID, err := s.OpenIDResolver.Resolve(request.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable,
			"cannot resolve openID:%v", err)
	}

	accountID, err := s.Mongo.ResolveAccountID(ctx, openID)
	if err != nil {
		s.Logger.Error("cannot resolve account id", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	token, err := s.TokenGenerator.GenerateToken(accountID.String(), s.TokenExpire)
	if err != nil {
		s.Logger.Error("cannot generate token")
		return nil, status.Error(codes.Internal, "")
	}

	return &authpb.LoginResponse{
		AccessToken: token,
		ExpiresIn:   int32(s.TokenExpire.Seconds()),
	}, nil
}
