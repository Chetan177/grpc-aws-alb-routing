package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"test/pkg/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return  &pb.HelloReply{}, fmt.Errorf("failed to get metadata")
	}

	log.Printf("Received Header: x-route-id = %v", md["x-route-id"])
	return &pb.HelloReply{Message: "Hello " + in.GetName() + " from Server-1"}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}