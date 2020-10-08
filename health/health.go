// Package health contains a gRPC and an HTTP server providing
// application health related information.
//
// Applications may use this package to provide liveness and readiness
// endpoints for Kubernetes probes.
//
// Application version information is also made available which should
// be set at build time with the `-X` linker flag. This helps identify
// exactly which version of the application is healthy or not.
package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sync/atomic"

	"github.com/anz-bank/pkg/health/pb"
	"google.golang.org/grpc"
)

// Undefined is the default value for the version strings. It exists to
// make it clear that the values have not been supplied at build time,
// as opposed to the empty string that is a common build-time
// miscalculated value.
const Undefined = "undefined"

// Variables served by the Version endpoint/method. These should be set at
// build time using the `-X` ldflag. e.g.
//     go build -ldflags='-X github.com/anz-bank/pkg/health.RepoURL="..."`
var (
	// RepoURL is the canonical repository source code URL.
	// e.g. https://github.com/anz-bank/pkg
	RepoURL = Undefined

	// CommitHash is the full git commit hash.
	// e.g. 1ee4e1f233caea38d6e331299f57dd86efb47361
	CommitHash = Undefined

	// BuildLogURL is the CI run URL.
	// e.g. https://github.com/anz-bank/pkg/actions/runs/84341844
	BuildLogURL = Undefined

	// ContainerTag is the canonical container image tag.
	// e.g. gcr.io/google-containers/hugo
	ContainerTag = Undefined

	// Semver is the semantic version compliant version.
	// e.g. v1.0.4
	Semver = Undefined

	// ScannerURLs is a JSON object containing URLs for additional code
	// scanner links.
	// e.g. { "codecov.io" : "https://codecov.io/..." }
	ScannerURLs string
)

// Errors defined and returned in this package.
var (
	// ErrInvalidSemver is a sentinel error returned when the Semver string
	// is not a valid semantic version string.
	ErrInvalidSemver = fmt.Errorf("invalid semver")
)

// The ReadyProvider interface allows for the ready state to be provided
// from the outside in a pull fashion. Typically there is not need to
// make use of this interface and just update the Ready state with
//
//    State.SetReady(true|false)
//
// However, if required, ReadyProvider can be specified as
//
//    State.ReadyProvider = localReadyProvider
//
// this turns State.SetReady into a no-op unless localReadyProvider is
// also a ReadySetter.
type ReadyProvider interface {
	IsReady() bool
}

// The ReadySetter interface is used in State.SetReady: If the
// State.ReadyProvider can be type asserted to be a ReadySetter it will
// be set. The default case wraps a boolean variable for setting and
// accessing.
type ReadySetter interface {
	SetReady(bool)
}

type readiness uint32

func (r *readiness) IsReady() bool {
	return atomic.LoadUint32((*uint32)(r)) == 1
}
func (r *readiness) SetReady(b bool) {
	var v uint32
	if b {
		v = 1
	}
	atomic.StoreUint32((*uint32)(r), v)
}

// State holds the state published by the health servers. Typically a
// single instance is shared amongst all servers.
type State struct {
	ReadyProvider
	Version *pb.VersionResponse
}

// NewState returns a State with the global version variables set in the
// Version field, and Ready as false. If any of the global version variables
// cannot be parsed, an error is returned.
func NewState() (*State, error) {
	v, err := newVersion()
	if err != nil {
		return nil, err
	}
	return &State{Version: v, ReadyProvider: new(readiness)}, nil
}

// SetReady sets the ready status served. The value can be changed as
// many times as is necessary over the lifetime of the application.
func (s *State) SetReady(ready bool) {
	if r, ok := s.ReadyProvider.(ReadySetter); ok {
		r.SetReady(ready)
	}
}

// SetReadyProvider sets the embedded ReadyProvider for state such that
// the ready value returned by state.IsReady() is ready from it.
func (s *State) SetReadyProvider(r ReadyProvider) {
	s.ReadyProvider = r
}

// Server is a server that can serve health data via gRPC and HTTP.
type Server struct {
	*State
	GRPC *GRPCServer
	HTTP *HTTPServer
}

// NewServer returns a health.Server implementing a gRPC and an HTTP server, to
// serve a common set of underlying health data. If any of the package-level
// version variables are invalid, an error is returned.
func NewServer() (*Server, error) {
	state, err := NewState()
	if err != nil {
		return nil, err
	}

	s := &Server{
		GRPC:  &GRPCServer{State: state},
		HTTP:  &HTTPServer{State: state},
		State: state,
	}
	return s, nil
}

// GRPCServer implements a gRPC interface for the Health service serving the
// anz.health.v1.Health service.
type GRPCServer struct {
	*State
	pb.UnimplementedHealthServer // embedded for forward compatible implementations
}

// NewGRPCServer returns a GRPCServer. If any of the package-level version
// variables are invalid, an error is returned.
func NewGRPCServer() (*GRPCServer, error) {
	state, err := NewState()
	if err != nil {
		return nil, err
	}
	return &GRPCServer{State: state}, nil
}

// RegisterWith registers the Health GRPCServer with the given grpc.Server.
func (g *GRPCServer) RegisterWith(s *grpc.Server) {
	pb.RegisterHealthServer(s, g)
}

// Alive implements the anz.health.v1.Health.Alive method returning an empty
// response. If the caller receives the response without error, it means that
// the application is alive.
func (g *GRPCServer) Alive(_ context.Context, _ *pb.AliveRequest) (*pb.AliveResponse, error) {
	return &pb.AliveResponse{}, nil
}

// Ready implements the anz.health.v1.Health.Ready method, returning a bool
// value indicating whether the application is ready to receive traffic. An
// application may become ready or not ready any number of times.
func (g *GRPCServer) Ready(_ context.Context, _ *pb.ReadyRequest) (*pb.ReadyResponse, error) {
	return &pb.ReadyResponse{Ready: g.State.IsReady()}, nil
}

// Version implements the anz.health.v1.Health.Version method, returning
// information to identify the running version of the application.
func (g *GRPCServer) Version(_ context.Context, _ *pb.VersionRequest) (*pb.VersionResponse, error) {
	return g.State.Version, nil
}

// HTTPServer implements an HTTP interface for the Health service at
// /healthz, /readyz and /version
type HTTPServer struct {
	*State
	mux *http.ServeMux
}

// NewHTTPServer returns an HTTPServer.
//
// HTTPServer implements http.Handler and serves HTTP responses on the
// following paths:
//
//   /healthz
//   /readyz
//   /version
//
// Use a custom http.Handler or http.ServerMux with HandleAlive,
// HandleReady and HandleVersion to serve on different URL paths.
//
// If any of the package-level version variables are invalid, an error
// is returned.
func NewHTTPServer() (*HTTPServer, error) {
	state, err := NewState()
	if err != nil {
		return nil, err
	}
	h := &HTTPServer{State: state, mux: http.NewServeMux()}
	h.RegisterWith(h.mux)
	return h, nil
}

// The Router interface allows HTTPServer.RegisterWith to work with any
// type of mux that implements the Handle method as defined on http.ServeMux.
type Router interface {
	Handle(path string, h http.Handler)
}

// RegisterWith registers our handlers for /healthz, /readyz and /version
// with the given Router. This allows for sharing the root namespace
// with other paths and not requiring that the caller register each handler
// individually.
func (h *HTTPServer) RegisterWith(r Router) {
	r.Handle("/healthz", requireGet(h.HandleAlive))
	r.Handle("/readyz", requireGet(h.HandleReady))
	r.Handle("/version", requireGet(h.HandleVersion))
}

// requireGet is http middleware that ensures that the request's method is GET.
// It has a slightly different signature to normal middleware, to suit our use case.
func requireGet(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			msg := fmt.Sprintf("%d method not allowed, use GET", http.StatusMethodNotAllowed)
			http.Error(w, msg, http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	})
}

// IsHealthEndpoint returns true if the request is for one of our health
// check endpoints (/healthz or /readyz). It is intended to be used with
// the OpenCensus ochttp plugin to not trace health checks.
//
// Use it in the ochttp.Handler:
//
//     ochttp.Handler{IsHealthEndpoint: health.IsHealthEndpoint, ...}
func IsHealthEndpoint(r *http.Request) bool {
	return r.URL.Path == "/healthz" || r.URL.Path == "/readyz"
}

// ServeHTTP implements http.Handler, handling GET requests for /healthz,
// /readyz and /version. Other methods on these paths will return
// 405 Method Not Allowed, and other paths will return 404 Not Found.
func (h *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// HandleAlive returns a 200 OK response. If the caller receives this, it means
// that the application is alive. Any other response should be treated as the
// application not being alive.
func (h *HTTPServer) HandleAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d ok\n", http.StatusOK)
}

// HandleReady returns a 200 OK response if the application is ready to receive
// traffic. It returns a 503 Service Unavailable response if it is not ready to
// receive traffic. An application may become ready or not ready any number of
// times.
func (h *HTTPServer) HandleReady(w http.ResponseWriter, r *http.Request) {
	if !h.State.IsReady() {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "%d service unavailable\n", http.StatusServiceUnavailable)
		return
	}
	fmt.Fprintf(w, "%d ok\n", http.StatusOK)
}

// HandleVersion returns a 200 OK response with a JSON body containing
// the application version information. It is the JSON-serialised form
// of the health.pb.VersionResponse struct.
func (h *HTTPServer) HandleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.MarshalIndent(h.State.Version, "", "  ")
	_, _ = w.Write(b)
}

func newVersion() (*pb.VersionResponse, error) {
	var scannerURLs map[string]string
	if ScannerURLs != "" {
		if err := json.Unmarshal([]byte(ScannerURLs), &scannerURLs); err != nil {
			return nil, fmt.Errorf("invalid ScannerURLs JSON: %w", err)
		}
	}
	version := &pb.VersionResponse{
		RepoUrl:      RepoURL,
		CommitHash:   CommitHash,
		BuildLogUrl:  BuildLogURL,
		ContainerTag: ContainerTag,
		Semver:       Semver,
		ScannerUrls:  scannerURLs,
	}

	if err := validateVersion(version); err != nil {
		return nil, err
	}
	return version, nil
}

var semverRe = regexp.MustCompile(`^` +
	`v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`)

func validateVersion(v *pb.VersionResponse) error {
	for s, u := range v.ScannerUrls {
		if _, err := url.ParseRequestURI(u); err != nil {
			return fmt.Errorf("invalid Scanner URL '%s': %w", s, err)
		}
	}
	if v.RepoUrl != Undefined {
		if _, err := url.ParseRequestURI(v.RepoUrl); err != nil {
			return fmt.Errorf("invalid RepoURL '%s': %w", v.RepoUrl, err)
		}
	}
	if v.BuildLogUrl != Undefined {
		if _, err := url.ParseRequestURI(v.BuildLogUrl); err != nil {
			return fmt.Errorf("invalid BuildLogURL '%s': %w", v.BuildLogUrl, err)
		}
	}
	if v.Semver != Undefined {
		if !semverRe.MatchString(v.Semver) {
			return fmt.Errorf("%w: %s", ErrInvalidSemver, v.Semver)
		}
	}
	return nil
}
