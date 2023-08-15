package main

import (
	"context"
	"github.com/lucaskatayama/learn-grpc/examples/simple-streaming/proto/ping"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"os"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("failed to dial", "err", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := ping.NewPingerClient(conn)

	resp, err := client.Ping(context.Background(), &ping.PingRequest{Uuid: "123"})
	if err != nil {
		slog.Error("failed to ping", "err", err)
		os.Exit(1)
	}

	done := make(chan int)
	go func() {
		for {
			d, err := resp.Recv()
			if err != nil {
				slog.Error("failed to received message", "err", err)
				return
			}
			slog.Info("message received", "created_at", d.CreatedAt)
		}
	}()

	<-done

}
