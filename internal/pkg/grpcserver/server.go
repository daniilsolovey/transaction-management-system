package grpcserver

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	Server *grpc.Server
	Port   int
}

func New(port int) *Server {
	return &Server{
		Server: grpc.NewServer(),
		Port:   port,
	}
}

func (s *Server) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		panic(err)
	}
	if err := s.Server.Serve(lis); err != nil {
		panic(err)
	}
}
