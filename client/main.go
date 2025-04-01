package main

import (
	"context"
	"crypto/tls"
	"flag"
	"io"
	"log"
	"time"

	pb "github.com/tuutoo/grpc-hello/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 命令行参数
	address := flag.String("addr", "localhost:50051", "gRPC server address")
	insecureConn := flag.Bool("insecure", false, "Use plaintext (non-TLS) connection")
	name := flag.String("name", "Jerry", "Name to greet")
	flag.Parse()

	// 设置连接选项
	var opts grpc.DialOption
	if *insecureConn {
		log.Println("Using insecure connection")
		opts = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		log.Println("Using TLS connection (InsecureSkipVerify=true)")
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true, // 用于测试，生产应验证证书
		}
		opts = grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))
	}

	conn, err := grpc.Dial(*address, opts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.SayHello(ctx, &pb.HelloRequest{Name: *name})
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
