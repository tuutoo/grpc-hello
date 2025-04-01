package main

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/tuutoo/grpc-hello/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(req *pb.HelloRequest, stream pb.Greeter_SayHelloServer) error {
	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("Hello %s - %d", req.Name, i+1)
		if err := stream.Send(&pb.HelloReply{Message: msg}); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &server{})
	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
