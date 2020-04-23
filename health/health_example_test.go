//nolint:gosec,errcheck
package health_test

import (
	"net"
	"net/http"

	"github.com/anz-bank/pkg/health"
	"google.golang.org/grpc"
)

func Example() {
	healthServer, _ := health.NewServer()
	go http.ListenAndServe(":8082", healthServer)

	grpcServer := grpc.NewServer()
	healthServer.RegisterWith(grpcServer)
	lis, _ := net.Listen("tcp", ":8080")
	go grpcServer.Serve(lis)

	// [expensive initialisation]

	healthServer.SetReady(true)
}
