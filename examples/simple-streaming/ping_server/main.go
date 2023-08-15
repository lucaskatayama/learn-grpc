package main

import (
	"github.com/lucaskatayama/learn-grpc/examples/simple-streaming/proto/ping"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
	"os"
	"time"
)

type server struct {
	ping.UnimplementedPingerServer
}

func (s server) Ping(request *ping.PingRequest, stream ping.Pinger_PingServer) error {

	for {
		select {
		case <-stream.Context().Done():
			slog.Info("connection closed")
			return status.Error(codes.Canceled, "request canceled")
		default:
			slog.Info("running ping")
			time.Sleep(1 * time.Second)
			if err := stream.SendMsg(&ping.PingReply{
				CreatedAt: time.Now().UTC().Format(time.RFC3339Nano),
			}); err != nil {
				slog.Error("failed", "err", err)
				return status.Error(codes.DataLoss, "failed to send message")
			}
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		slog.Error("failed to listen to localhost:50051", "err", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	ping.RegisterPingerServer(s, &server{})
	slog.Info("server listening", "port", lis.Addr())
	if err := s.Serve(lis); err != nil {
		slog.Error("failed to serve", "error", err)
		os.Exit(1)
	}
}
