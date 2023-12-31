package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"strings"

	// importing generated stubs
	gen "grpc-gateway-demo/gen/go/hello"
)

func getApiVersion(ctx context.Context) string {
	apiVersion := ""
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if xApiVersion, ok := md["x-api-version"]; ok {
			apiVersion = strings.Join(xApiVersion, ",")
		}
	}
	return apiVersion
}

func setPagingAndSortingInfo() metadata.MD {

	md := metadata.Pairs(
		"x-page-no", "2",
		"x-page-size", "10",
		"x-total-elements", "50",
		"x-total-pages", "5")
	return md

}

// GreeterServerImpl will implement the service defined in protocol buffer definitions
type GreeterServerImpl struct {
	gen.UnimplementedGreeterServer
}

// SayHello is the implementation of RPC call defined in protocol definitions.
// This will take HelloRequest message and return HelloReply
func (g *GreeterServerImpl) SayHello(ctx context.Context, request *gen.HelloRequest) (*gen.HelloReply, error) {
	version := getApiVersion(ctx)
	err := grpc.SendHeader(ctx, setPagingAndSortingInfo())
	if err != nil {
		return nil, err
	}
	return &gen.HelloReply{
		Message:    fmt.Sprintf("hello %s", request.Name),
		ApiVersion: version,
	}, nil
}
func main() {
	// create new gRPC server
	server := grpc.NewServer()
	// register the GreeterServerImpl on the gRPC server
	gen.RegisterGreeterServer(server, &GreeterServerImpl{})

	log.Println("Starting gRPC server on http://0.0.0.0:8080")
	// start listening on port :8080 for a tcp connection
	if l, err := net.Listen("tcp", ":8080"); err != nil {
		log.Fatal("error in listening on port :8080", err)
	} else {
		// the gRPC server
		if err := server.Serve(l); err != nil {
			log.Fatal("unable to start server", err)
		}
	}

}
