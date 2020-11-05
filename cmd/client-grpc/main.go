package main

import (
	"context"
	"flag"
	"log"
	"time"

	v1 "github.com/aaabhilash97/op/pkg/api/v1"
	"google.golang.org/grpc"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

func main() {
	// get configuration
	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := v1.NewOpServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call Create
	req1 := v1.NewUserRegistrationRequest{}
	res1, err := c.NewUserRegistration(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", res1)

}
