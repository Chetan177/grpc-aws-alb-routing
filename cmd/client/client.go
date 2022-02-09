package main

import (
	"context"
	"flag"
	"log"
	"time"

	"test/pkg/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "GRPC-TEST-....amazonaws.com:443", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	var opts []grpc.DialOption
	// change /test/mykey.pem with your cetificate file
	caFile:= "/test/mykey.pem" 
	// change test.example.com with your domain
	creds, err := credentials.NewClientTLSFromFile(caFile, "test.example.com")
	if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))

	conn, err := grpc.Dial(*addr, opts...)
	if err != nil {
			log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	log.Println("sending request to Server-1 with flag blue")
	md := metadata.New(map[string]string{"x-route-id": "blue"})
	ctx = metadata.NewOutgoingContext(ctx, md)
	
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	log.Println("sending request to Server-2 with flag green")
	md = metadata.New(map[string]string{"x-route-id": "green"})
	ctx = metadata.NewOutgoingContext(ctx, md)
	
	r, err = c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
