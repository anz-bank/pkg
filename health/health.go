package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	pb "github.com/anz-bank/pkg/health/healthpb"
	"google.golang.org/grpc"
)

const Undefined = "undefined"

var (
	RepoURL      = Undefined
	CommitHash   = Undefined
	BuildLogURL  = Undefined
	ContainerTag = Undefined
	Semver       = Undefined
	ScannerURLs  string

	ErrInvalidSemver = fmt.Errorf("invalid semver")
)

// --- Server ----------------------------------------------------------------

// Server ...
type Server struct {
	GRPC *GRPCServer
	HTTP *HTTPServer
	data *serverData
}

// NewServer ...
func NewServer() (*Server, error) {
	data, err := newServerData()
	if err != nil {
		return nil, err
	}

	s := &Server{
		GRPC: NewGRPCServer(data),
		HTTP: NewHTTPServer(data),
		data: data,
	}
	return s, nil
}

func (s *Server) SetReady(ready bool) {
	s.data.SetReady(ready)
}

// --- GRPCServer ------------------------------------------------------------

// GRPCServer implements a gRPC interface for the Health service serving the
// health data supplied a Service.
type GRPCServer struct {
	pb.UnimplementedHealthServer
	data *serverData
}

// NewGRPCServer returns a GRPCServer that serves the data provided by the
// Service interface provided.
func NewGRPCServer(data *serverData) *GRPCServer {
	return &GRPCServer{data: data}
}

// RegisterWith registers the Health GRPCServer with the given grpc.Server.
func (g *GRPCServer) RegisterWith(s *grpc.Server) {
	pb.RegisterHealthServer(s, g)
}

// Alive ...
func (g *GRPCServer) Alive(_ context.Context, _ *pb.AliveRequest) (*pb.AliveResponse, error) {
	return &pb.AliveResponse{}, nil
}

// Ready ...
func (g *GRPCServer) Ready(_ context.Context, _ *pb.ReadyRequest) (*pb.ReadyResponse, error) {
	return &pb.ReadyResponse{Ready: g.data.Ready()}, nil
}

// Version ...
func (g *GRPCServer) Version(_ context.Context, _ *pb.VersionRequest) (*pb.VersionResponse, error) {
	return g.data.Version(), nil
}

// --- HTTPServer ------------------------------------------------------------

// HTTPServer implements an HTTP service for the Health service serving the
// health data supplied by Service.
type HTTPServer struct {
	data *serverData
}

// NewHTTPServer
func NewHTTPServer(data *serverData) *HTTPServer {
	return &HTTPServer{data: data}
}

// ServeHTTP ...
func (h *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		msg := fmt.Sprintf("%d method not allowed, use GET", http.StatusMethodNotAllowed)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}
	switch r.URL.Path {
	case "/healthz":
		h.HandleAlive(w, r)
	case "/readyz":
		h.HandleReady(w, r)
	case "/version":
		h.HandleVersion(w, r)
	default:
		http.NotFound(w, r)
	}
}

// HandleAlive ...
func (h *HTTPServer) HandleAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d ok\n", http.StatusOK)
}

// HandleReady ...
func (h *HTTPServer) HandleReady(w http.ResponseWriter, r *http.Request) {
	if !h.data.Ready() {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "%d service unavailable\n", http.StatusServiceUnavailable)
		return
	}
	fmt.Fprintf(w, "%d ok\n", http.StatusOK)
}

// HandleVersion ...
func (h *HTTPServer) HandleVersion(w http.ResponseWriter, r *http.Request) {
	b, _ := json.MarshalIndent(h.data.Version(), "", "  ")
	_, _ = w.Write(b)
}

// --- serverData -------------------------------------------------------------

type serverData struct {
	ready   bool
	version *pb.VersionResponse
}

func newServerData() (*serverData, error) {
	v, err := NewVersion()
	if err != nil {
		return nil, err
	}

	return &serverData{version: v}, nil
}

// SetReady ...
func (sd *serverData) SetReady(ready bool) {
	sd.ready = ready
}

// Ready ...
func (sd *serverData) Ready() bool {
	return sd.ready
}

// Version ...
func (sd *serverData) Version() *pb.VersionResponse {
	return sd.version
}

// NewVersion
func NewVersion() (*pb.VersionResponse, error) {
	var scannerURLs map[string]string
	if ScannerURLs != "" {
		if err := json.Unmarshal([]byte(ScannerURLs), &scannerURLs); err != nil {
			return nil, err
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
			return fmt.Errorf("%s: %w", s, err)
		}
	}
	if v.RepoUrl != Undefined {
		if _, err := url.ParseRequestURI(v.RepoUrl); err != nil {
			return err
		}
	}
	if v.BuildLogUrl != Undefined {
		if _, err := url.ParseRequestURI(v.BuildLogUrl); err != nil {
			return err
		}
	}
	if v.Semver != Undefined {
		if !semverRe.MatchString(v.Semver) {
			return fmt.Errorf("%w: %s", ErrInvalidSemver, v.Semver)
		}
	}
	return nil
}
