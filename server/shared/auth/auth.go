package token

import (
	"context"
	"coolcar/shared/id"
	"fmt"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"os"
	"strings"
)

const (
	authorizationHeader = "authorization"
	tokenPrefix         = "Bearer "
)

type interceptor struct {
	verifier tokenVerifier
}

func Interceptor(publicKeyFile string) (grpc.UnaryServerInterceptor, error) {
	file, err := os.Open(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("can not open public key file:%v", err)
	}

	all, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("can not read public key file:%v", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(all)
	if err != nil {
		return nil, fmt.Errorf("can not parse public key:%v")
	}

	i := &interceptor{
		verifier: &JWTVerifier{
			PublicKey: publicKey,
		},
	}
	return i.HandleReq, nil
}

type tokenVerifier interface {
	Verify(token string) (string, error)
}

func (i *interceptor) HandleReq(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	tkn, err := tokenFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "")
	}

	aid, err := i.verifier.Verify(tkn)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token is not valid %v", err)
	}

	return handler(ContextWithAccountID(ctx, id.AccountID(aid)), req)
}

func tokenFromContext(c context.Context) (string, error) {
	unauthenticated := status.Error(codes.Unauthenticated, "")

	m, ok := metadata.FromIncomingContext(c)

	if !ok {
		return "", unauthenticated
	}

	tkn := ""
	for _, v := range m[authorizationHeader] {
		if strings.HasPrefix(v, tokenPrefix) {
			tkn = v[len(tokenPrefix):]
		}

		if tkn == "" {
			return "", unauthenticated
		}
	}
	return tkn, nil
}

type accountIDKey struct{}

func ContextWithAccountID(c context.Context, aid id.AccountID) context.Context {
	return context.WithValue(c, accountIDKey{}, aid)
}

func AccountIDFromContext(c context.Context) (id.AccountID, error) {
	v := c.Value(accountIDKey{})

	aid, ok := v.(id.AccountID)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "")
	}

	return aid, nil
}
