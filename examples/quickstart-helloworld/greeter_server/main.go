package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/lucaskatayama/learn-grpc/examples/quickstart-helloworld/helloworld"
	"google.golang.org/grpc"
	log "log/slog"
	"net"
	"os"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Info("Received", "message", request.GetName())
	return &pb.HelloReply{Message: "Hello " + request.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Error("failed to listen", "error", err)
		os.Exit(1)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Info("server listening", "port", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Error("failed to serve", "error", err)
		os.Exit(1)
	}
}
