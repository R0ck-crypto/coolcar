package blob

import (
	"context"
	blobbpb "coolcar/blob/api/gen/v1"
	"coolcar/blob/dao"
	"coolcar/shared/id"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Storage
type Storage interface {
	SignURL(ctx context.Context, method, path string, timeout time.Duration) (string, error)
	Get(ctx context.Context, path string) (io.ReadCloser, error)
}

type Service struct {
	Storage Storage
	Mongo   *dao.Mongo
	Logger  *zap.Logger
	blobbpb.UnimplementedBlobServiceServer
}

func (s Service) CreateBlob(ctx context.Context, request *blobbpb.CreateBlobRequest) (*blobbpb.CreateBlobResponse, error) {
	aid := id.AccountID(request.AccountId)

	// 在数据库创建记录
	br, err := s.Mongo.CreateBlob(ctx, aid)
	if err != nil {
		s.Logger.Error("cannot create blob", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}
	// 生成预签名的上传URL
	url, err := s.Storage.SignURL(ctx, http.MethodPut, br.Path, secToDuration(request.UploadUrlTimeoutSec))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "cannot sign url : %v", err)
	}

	return &blobbpb.CreateBlobResponse{
		Id:        br.ID.Hex(),
		UploadUrl: url,
	}, nil
}

func (s Service) GetBlob(ctx context.Context, request *blobbpb.GetBlobRequest) (*blobbpb.GetBlobResponse, error) {
	br, err := s.getBlobRecord(ctx, id.BlobID(request.Id))
	if err != nil {
		return nil, err
	}
	get, err := s.Storage.Get(ctx, br.Path)
	if get != nil {
		defer get.Close()
	}

	if err != nil {
		return nil, status.Errorf(codes.Aborted, "cannot get stroage")
	}

	b, err := ioutil.ReadAll(get)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "cannot read from response: %v", err)
	}

	return &blobbpb.GetBlobResponse{Data: b}, nil

}

func (s Service) GetBlobURL(ctx context.Context, request *blobbpb.GetBlobURLRequest) (*blobbpb.GetBlobURLResponse, error) {
	br, err := s.getBlobRecord(ctx, id.BlobID(request.Id))
	if err != nil {
		return nil, err
	}

	// 生成预签名的下载URL
	url, err := s.Storage.SignURL(ctx, http.MethodGet, br.Path, secToDuration(request.TimeoutSec))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "cannot signed url")
	}

	return &blobbpb.GetBlobURLResponse{
		Url: url,
	}, nil
}

func (s Service) getBlobRecord(ctx context.Context, bid id.BlobID) (*dao.BlobRecord, error) {
	blob, err := s.Mongo.GetBlob(ctx, bid)
	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.NotFound, "")
	}
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return blob, nil
}

func secToDuration(sec int32) time.Duration {
	return time.Duration(sec) * time.Second
}
