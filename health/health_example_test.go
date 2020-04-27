//nolint:gosec,errcheck
package health_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"

	"github.com/anz-bank/pkg/health"
	"github.com/anz-bank/pkg/health/pb"
	"google.golang.org/grpc"
)

func Example() {
	server, _ := health.NewServer()
	go http.ListenAndServe(":8082", server.HTTP)

	gs := grpc.NewServer()
	server.GRPC.RegisterWith(gs)
	lis, _ := net.Listen("tcp", ":8080")
	go gs.Serve(lis)

	// [run expensive initialisation]

	server.SetReady(true)
}

func ExampleNewHTTPServer() {
	server, err := health.NewHTTPServer()
	if err != nil {
		log.Fatal(err)
	}
	// go http.ListenAndServe(":8082", server)

	r := httptest.NewRequest("GET", "/readyz", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, r)
	fmt.Print(w.Body.String())

	// [run expensive initialisation]

	server.SetReady(true)

	w = httptest.NewRecorder()
	server.ServeHTTP(w, r)
	fmt.Print(w.Body.String())
	// output: 503 service unavailable
	// 200 ok
}

func ExampleNewGRPCServer() {
	server, err := health.NewGRPCServer()
	if err != nil {
		log.Fatal(err)
	}
	// go grpcListenAndServe(":8082", server)

	ctx := context.Background()
	resp, err := server.Ready(ctx, &pb.ReadyRequest{})
	fmt.Println("err:", err, "ready:", resp.Ready)

	// [run expensive initialisation]

	server.SetReady(true)

	resp, err = server.Ready(ctx, &pb.ReadyRequest{})
	fmt.Println("err:", err, "ready:", resp.Ready)
	// output: err: <nil> ready: false
	// err: <nil> ready: true
}
