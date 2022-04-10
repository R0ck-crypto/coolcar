package server

import (
	token "coolcar/shared/auth"
	zap "go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GRPConfig struct {
	Name              string
	Addr              string
	AuthPublickeyfile string
	RegisterFunc      func(server *grpc.Server)
	Logger            *zap.Logger
}

func RunGRPCServer(c *GRPConfig) error {
	nameFelds := zap.String("name", c.Name)
	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		c.Logger.Fatal("cannot listen", nameFelds, zap.Error(err))
	}

	var opts []grpc.ServerOption
	if c.AuthPublickeyfile != "" {
		in, err := token.Interceptor("shared/auth/public.key")
		if err != nil {
			c.Logger.Fatal("cannot create auth interceptor", nameFelds, zap.Error(err))
		}
		opts = append(opts, grpc.UnaryInterceptor(in))
	}

	s := grpc.NewServer(opts...)

	c.RegisterFunc(s)

	c.Logger.Info("server started", nameFelds, zap.String("addr", c.Addr))
	return s.Serve(lis)

}
