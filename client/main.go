package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/tuutoo/grpc-hello/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("192.168.66.203:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Jerry"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error receiving: %v", err)
		}
		log.Printf("Received: %s", res.Message)
	}
}
