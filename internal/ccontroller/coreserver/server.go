package coreserver

import (
	"context"
	"crypto/tls"
	"github.com/eskpil/sunlight/internal/ccontroller/mycontext"
	"github.com/eskpil/sunlight/pkg/api/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
)

type Server struct {
	// TODO: Figure out how to remove this
	core.UnimplementedCoreServer

	tlsc     tls.Config
	listener net.Listener

	c *mycontext.Context
}

func NewServer() (*Server, error) {
	server := new(Server)

	cert, err := tls.LoadX509KeyPair(os.ExpandEnv("$SUNLIGHT_PKI_DIR/certs/cert.crt"), os.ExpandEnv("$SUNLIGHT_PKI_DIR/keys/key.key"))
	if err != nil {
		return nil, err
	}
	server.tlsc = tls.Config{Certificates: []tls.Certificate{cert}}

	listener, err := net.Listen("tcp", "127.0.0.1:1900")
	if err != nil {
		return nil, err
	}
	server.listener = listener

	return server, nil
}

func (s *Server) UpdateResource(ctx context.Context, resource *core.Resource) (*core.UpdateResourceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) CreateResource(ctx context.Context, resource *core.Resource) (*core.CreateResourceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) GetResource(ctx context.Context, request *core.GetResourceRequest) (*core.GetResourceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) mustEmbedUnimplementedCoreServer() {
	//TODO implement me
	panic("implement me")
}

func (s *Server) Start() error {
	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(&s.tlsc)),
	)

	core.RegisterCoreServer(grpcServer, s)

	return grpcServer.Serve(s.listener)
}
