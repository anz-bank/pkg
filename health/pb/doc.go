//go:generate prototool lint health.proto
//go:generate prototool format -w
//go:generate protoc --go_out=plugins=grpc,paths=source_relative:. health.proto
//go:generate goimports -w health.pb.go

// Package pb provides protoc generated types and methods for
// Health service with Alive, Ready and Version methods.
package pb
