package main

import (
	"context"
	"fmt"
	"github.com/lucaskatayama/learn-grpc/examples/simple-streaming/proto/ping"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"os"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("failed to dial", "err", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := ping.NewPingerClient(conn)

	// server streaming
	slog.Info("-------------------- Server Streaming -----------------------")
	resp, err := client.ServerStream(context.Background(), &ping.PingRequest{Uuid: "123"})
	if err != nil {
		slog.Error("failed to ping", "err", err)
		os.Exit(1)
	}

	count := 0
	for {
		if count > 5 {
			break
		}
		d, err := resp.Recv()
		if err != nil {
			slog.Error("failed to received message", "err", err)
			break
		}
		slog.Info("message received", "created_at", d.CreatedAt)
		count++
	}

	// client streaming
	slog.Info("-------------------- Client Streaming -----------------------")
	count = 0
	clientStream, err := client.ClientStream(context.Background())
	if err != nil {
		slog.Error("creating client streaming", "err", err)
		os.Exit(1)
	}
	for {
		if count > 5 {
			break
		}
		time.Sleep(1 * time.Second)
		if err := clientStream.Send(&ping.PingRequest{Uuid: fmt.Sprintf("ID#%d", count)}); err != nil {
			slog.Error("sending stream message", "err", err)
			break
		}
		count++
	}

	// bidi streaming
	slog.Info("-------------------- Bidi Streaming -----------------------")
	count = 0
	bidiStream, err := client.BidiStream(context.Background())
	if err != nil {
		slog.Error("creating client streaming", "err", err)
		os.Exit(1)
	}

	go func() {
		for {
			d, err := bidiStream.Recv()
			if err != nil {
				slog.Error("reading message", "err", err)
				return
			}
			slog.Info("reaceived", "msg", d.CreatedAt)
		}
	}()
	for {
		if count > 5 {
			break
		}
		time.Sleep(1 * time.Second)
		if err := bidiStream.Send(&ping.PingRequest{Uuid: fmt.Sprintf("ID#%d", count)}); err != nil {
			slog.Error("sending stream message", "err", err)
			os.Exit(1)
		}
		count++
	}

}
