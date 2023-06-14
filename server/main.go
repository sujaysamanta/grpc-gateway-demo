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

// GreeterServerImpl will implement the service defined in protocol buffer definitions
type GreeterServerImpl struct {
	gen.UnimplementedGreeterServer
}

// SayHello is the implementation of RPC call defined in protocol definitions.
// This will take HelloRequest message and return HelloReply
func (g *GreeterServerImpl) SayHello(ctx context.Context, request *gen.HelloRequest) (*gen.HelloReply, error) {
	version := getApiVersion(ctx)
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
	// start listening on port :8080 for a tcp connection
	if l, err := net.Listen("tcp", ":8080"); err != nil {
		log.Fatal("error in listening on port :8080", err)
	} else {
		// the gRPC server
		if err := server.Serve(l); err != nil {
			log.Fatal("unable to start server", err)
		}
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8081")
}
