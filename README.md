# pkg

[![Go-Linux](https://github.com/anz-bank/pkg/workflows/Go-Linux/badge.svg)](https://github.com/anz-bank/pkg/actions?query=workflow%3AGo-Linux+branch%3Amaster)
[![Godoc](https://img.shields.io/badge/godoc-ref-blue)](https://pkg.go.dev/github.com/anz-bank/pkg)
[![Slack chat](https://img.shields.io/badge/slack-anzoss-795679?logo=slack)](https://anzoss.slack.com/app_redirect?channel=pkg)

Common ANZ Go Packages

### Development

-   Pre-requisites: [go](https://golang.org/doc/go1.18),
    [golangci-lint](https://github.com/golangci/golangci-lint/releases/tag/v1.48.0),
    GNU make
-   Build with `make`
-   View build options with `make help`

On OSX, after installing go [1.18](https://golang.org/doc/install) run

    brew install golangci/tap/golangci-lint make

### Working with protos

All generated code derived from Protobuf and gRPC definitions is
committed to this repo, however if you need to regenerate the Go code or
work with gRPC, install the following tools:

-   proto3 and gRPC
    -   https://github.com/protocolbuffers/protobuf/releases
    -   https://github.com/golang/protobuf
    -   https://github.com/grpc/grpc
-   [`prototool`](https://github.com/uber/prototool/blob/dev/docs/install.md)
-   [`goimports`](https://godoc.org/golang.org/x/tools/cmd/goimports)
-   [`gprcurl`](https://github.com/fullstorydev/grpcurl)

On OSX run

    (cd /tmp; go get -u golang.org/x/tools/goimports)
    brew install grpcurl protoc-gen-go grpc prototool

After the initial installation run

    make generate
