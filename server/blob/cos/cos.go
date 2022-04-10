package cos

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Service struct {
	client *cos.Client
	secID  string
	secKey string
}

func (s Service) SignURL(ctx context.Context, method, path string, timeout time.Duration) (string, error) {
	u, err := s.client.Object.GetPresignedURL(ctx, method, path, s.secID,
		s.secKey, timeout, nil)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

func (s Service) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	res, err := s.client.Object.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var b io.ReadCloser
	if res != nil {
		b = res.Body
	}

	if err != nil {
		return b, err
	}

	if res.StatusCode >= 400 {
		return b, fmt.Errorf("got err response :%+v", res)
	}

	return b, nil

}

func NewService(addr string, secID string, secKey string) (*Service, error) {

	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse addr : %v", addr)
	}

	b := &cos.BaseURL{BucketURL: u, ServiceURL: u}
	// 1.永久密钥

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secID,  // 替换为用户的 SecretId，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
			SecretKey: secKey, // 替换为用户的 SecretKey，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
		},
	})

	return &Service{
		client: client,
		secID:  secID,
		secKey: secKey,
	}, nil

}
