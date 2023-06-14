package main

import (
	"context"
	"github.com/felixge/httpsnoop"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	gen "grpc-gateway-demo/gen/go/hello"
	"log"
	"net"
	"net/http"
)

var DefaultApiVersion = "1.0.0"

func withLogger(handler http.Handler) http.Handler {
	// the creation a handler
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// pass the handler to https-noop to get http status and latency
		metrics := httpsnoop.CaptureMetrics(handler, writer, request)
		// printing exacted data
		log.Printf("http[%d]-- %s -- %s\n", metrics.Code, metrics.Duration, request.URL.Path)
	})
}

/*
func CustomMatcher(key string) (string, bool) {
	switch key {
	case "X-User-Id":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
*/

func defaultToVersion(ctx context.Context, request *http.Request) metadata.MD {
	header := request.Header.Get("x-api-version")
	var md metadata.MD = nil
	if header == "" {
		md = metadata.Pairs("x-api-version", DefaultApiVersion)
	} else {
		md = metadata.Pairs("x-api-version", header)
	}
	return md
}

func main() {
	// creating mux for gRPC gateway. This will multiplex or route request different gRPC service
	mux := runtime.NewServeMux(
		//runtime.WithIncomingHeaderMatcher(CustomMatcher),
		runtime.WithMetadata(defaultToVersion),
	)
	// setting up a dail up for gRPC service by specifying endpoint/target url
	err := gen.RegisterGreeterHandlerFromEndpoint(context.Background(), mux, "localhost:8080", []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatal(err)
	}

	// Creating a normal HTTP server
	server := http.Server{
		Handler: withLogger(mux),
	}
	// creating a listener for server
	listener, err := net.Listen("tcp", ":8081")

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Staring to serve gRPC-Gateway on http://0.0.0.0:8081")

	// start server
	err = server.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
