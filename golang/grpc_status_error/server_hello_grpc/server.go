package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"sync/atomic"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/reflection"
)

var counter uint64 = 0

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("pb.HelloRequest", in.Name)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Printf("get metadata error")
	}
	for key, val := range md {
		fmt.Printf("%v:%v\n", key, val)
	}

	x := atomic.AddUint64(&counter, 1) % 3
	switch x {
	case 0:
		return &pb.HelloReply{Message: "Hello " + in.Name}, nil
	case 1:
		return nil, status.Error(codes.DataLoss, "--DataLoss--")
	default:
		return nil, status.Error(codes.Unauthenticated, "--Unauthenticated--")
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("say_hello_grpc starting...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
