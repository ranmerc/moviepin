package grpcserver

import (
	"net"
	"token-management-service/tokenpb"

	"google.golang.org/grpc"
)

// Server implements gRPC server for token management service
type Server struct {
	port          string
	serverHandler tokenpb.TokenServiceServer
}

// NewServer returns new instance of tms gRPC server
func NewServer(port string, serverHandler tokenpb.TokenServiceServer) *Server {
	return &Server{
		port:          port,
		serverHandler: serverHandler,
	}
}

// ListenAndServe servers the gRPC server on specified port
func (s *Server) ListenAndServe() error {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	tokenpb.RegisterTokenServiceServer(grpcServer, s.serverHandler)

	if err := grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
