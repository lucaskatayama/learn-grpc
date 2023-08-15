package main

import (
	"fmt"
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

func (s server) ServerStream(request *ping.PingRequest, stream ping.Pinger_ServerStreamServer) error {
	count := 0
	for {
		if count > 5 {
			return nil
		}
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
				return status.Error(codes.Canceled, "failed to send message")
			}
		}
		count++
	}
}

func (s server) ClientStream(stream ping.Pinger_ClientStreamServer) error {
	for {
		d, err := stream.Recv()
		if err != nil {
			slog.Error("receiving message", "err", err)
			return status.Error(codes.Canceled, "stream canceled")
		}
		slog.Info("message", "msg", d.Uuid)
	}
}

func (s server) BidiStream(stream ping.Pinger_BidiStreamServer) error {
	for {
		d, err := stream.Recv()
		if err != nil {
			slog.Error("receiving message", "err", err)
			return status.Error(codes.Canceled, "stream canceled")
		}
		slog.Info("message", "msg", d.Uuid)
		if err := stream.Send(&ping.PingReply{
			CreatedAt: fmt.Sprintf("[%s] -> %s", d.Uuid, time.Now().Format(time.RFC3339Nano)),
		}); err != nil {
			slog.Error("receiving message", "err", err)
			return status.Error(codes.Canceled, "stream canceled")
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
