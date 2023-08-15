package main

import (
	"context"
	"flag"
	pb "github.com/lucaskatayama/learn-grpc/examples/quickstart-helloworld/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	log "log/slog"
	"os"
	"time"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("did not connect", "error", err)
		os.Exit(1)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Error("could not greet", "error", err)
		os.Exit(1)
	}
	log.Info("Greeting", "message", r.GetMessage())
}
