package health

import (
	"sync"

	"google.golang.org/grpc"
)

var (
	// DefaultServer is the Server instance on which the package-level functions
	// operate. Most servers need only a single health server so this provides
	// a convenient definition of it that is available everywhere.
	DefaultServer *Server
	defaultState  = State{ReadyProvider: new(readiness)}

	serverInit sync.Once
)

// RegisterWithGRPC registers the default server
// health.DefaultServer.GRPC with the given grpc Server to make the
// health service available at "anz.health.v1.Health". This
// RegisterWithGRPC function returns an error when the Version
// information is invalid.
func RegisterWithGRPC(s *grpc.Server) error {
	var err error
	serverInit.Do(func() { err = newDefaultServer() })
	if err != nil {
		return err
	}
	DefaultServer.GRPC.RegisterWith(s)
	return nil
}

// RegisterWithHTTP registers the default server
// health.DefaultServer.HTTP with the given Router, e.g. a
// http.ServeMux, to make the health service endpoints available at
// /healthz, /readyz and /version. This RegisterWithHTTP function
// returns an error when the Version information is invalid.
func RegisterWithHTTP(r Router) error {
	var err error
	serverInit.Do(func() { err = newDefaultServer() })
	if err != nil {
		return err
	}
	DefaultServer.HTTP.RegisterWith(r)
	return nil
}

// SetReady sets the ready status served by the DefaultServer. The value
// can be changed as many times as is necessary over the lifetime of the
// application. It is valid to call SetReady before the DefaultServer
// has been registered, however a invalid version information will only
// be detected when RegisterWithHTTP or RegisterWithGRPC is called.
func SetReady(ready bool) {
	defaultState.SetReady(ready)
}

// SetReadyProvider sets the ReadyProvider for the DefaultServer.
func SetReadyProvider(r ReadyProvider) {
	defaultState.SetReadyProvider(r)
}

func newDefaultServer() error {
	v, err := newVersion()
	if err != nil {
		return err
	}
	defaultState.Version = v

	DefaultServer = &Server{
		GRPC:  &GRPCServer{State: &defaultState},
		HTTP:  &HTTPServer{State: &defaultState},
		State: &defaultState,
	}
	return nil
}
