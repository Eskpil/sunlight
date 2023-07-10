package adoption

import (
	"context"
	"crypto/tls"
	"github.com/eskpil/sunlight/pkg/api/adoption"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
)

type Server struct {
	adoption.UnimplementedAdoptionServer

	tlsc     tls.Config
	listener net.Listener
}

func NewServer() (*Server, error) {
	server := new(Server)

	cert, err := tls.LoadX509KeyPair(os.ExpandEnv("$SUNLIGHT_PKI_DIR/certs/cert.crt"), os.ExpandEnv("$SUNLIGHT_PKI_DIR/keys/key.key"))
	if err != nil {
		return nil, err
	}
	server.tlsc = tls.Config{Certificates: []tls.Certificate{cert}}

	listener, err := net.Listen("tcp", "127.0.0.1:2001")
	if err != nil {
		return nil, err
	}
	server.listener = listener

	return server, nil
}

func (s *Server) NeedsAdoption(ctx context.Context, request *adoption.NeedsAdoptionRequest) (*adoption.NeedsAdoptionResponse, error) {
	//TODO implement me

	response := &adoption.NeedsAdoptionResponse{
		Verdict: true,
		Hints: &adoption.AdoptionHints{
			CommonName:         "workstation1.",
			Locality:           "Trondheim",
			Country:            "Norway",
			Province:           "Trondelag",
			Organization:       "Eskpil",
			OrganizationalUnit: "workstations",
			Email:              "admin@workstation.local",
		},
		Requirements: []string{"tpm2"},
		Error:        nil,
	}

	return response, nil
}

func (s *Server) Adopt(ctx context.Context, request *adoption.AdoptRequest) (*adoption.AdoptResponse, error) {
	response := &adoption.AdoptResponse{
		Error:       nil,
		Verdict:     true,
		Certificate: "",
	}

	return response, nil
}

func (s *Server) Start() error {
	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(&s.tlsc)),
	)

	adoption.RegisterAdoptionServer(grpcServer, s)

	return grpcServer.Serve(s.listener)
}
