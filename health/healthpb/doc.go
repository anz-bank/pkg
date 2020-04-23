//go:generate prototool lint health.proto
//go:generate protoc --go_out=plugins=grpc:. health.proto
//
// package health provides liveliness, readiness and version endpoints.
package healthpb
