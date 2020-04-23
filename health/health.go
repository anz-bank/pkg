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

// Server ...
type Server struct {
	grpcServer *GRPCServer
	httpServer *HTTPServer
}

// NewServer ...
func NewServer() (*Server, error) {
	g, err := NewGRPCServer()
	if err != nil {
		return nil, err
	}
	h, err := NewHTTPServer()
	if err != nil {
		return nil, err
	}
	return &Server{grpcServer: g, httpServer: h}, nil
}

// SetReady ...
func (s *Server) SetReady(ready bool) {
	s.grpcServer.SetReady(ready)
	s.httpServer.SetReady(ready)
}

// RegisterWith ...
func (s *Server) RegisterWith(g *grpc.Server) {
	s.grpcServer.RegisterWith(g)
}

// RegisterWith ...
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.httpServer.ServeHTTP(w, r)
}

// Server ...
type GRPCServer struct {
	pb.UnimplementedHealthServer

	ready   bool
	version *pb.VersionResponse
}

func NewGRPCServer() (*GRPCServer, error) {
	version, err := NewVersion()
	if err != nil {
		return nil, err
	}
	return &GRPCServer{version: version}, nil
}

// Alive ...
func (*GRPCServer) Alive(_ context.Context, _ *pb.AliveRequest) (*pb.AliveResponse, error) {
	return &pb.AliveResponse{}, nil
}

// Ready ...
func (g *GRPCServer) Ready(_ context.Context, _ *pb.ReadyRequest) (*pb.ReadyResponse, error) {
	return &pb.ReadyResponse{Ready: g.ready}, nil
}

// Version ...
func (g *GRPCServer) Version(_ context.Context, _ *pb.VersionRequest) (*pb.VersionResponse, error) {
	return g.version, nil
}

// SetReady ...
func (g *GRPCServer) SetReady(ready bool) {
	g.ready = ready
}

// RegisterWith ...
func (g *GRPCServer) RegisterWith(s *grpc.Server) {
	pb.RegisterHealthServer(s, g)
}

// HTTPServer ...
type HTTPServer struct {
	ready   bool
	version *pb.VersionResponse
}

// NewHTTPServer
func NewHTTPServer() (*HTTPServer, error) {
	version, err := NewVersion()
	if err != nil {
		return nil, err
	}
	return &HTTPServer{version: version}, nil
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
func (*HTTPServer) HandleAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d ok\n", http.StatusOK)
}

// HandleReady ...
func (h *HTTPServer) HandleReady(w http.ResponseWriter, r *http.Request) {
	if !h.ready {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "%d service unavailable\n", http.StatusServiceUnavailable)
		return
	}
	fmt.Fprintf(w, "%d ok\n", http.StatusOK)
}

// HandleVersion ...
func (h *HTTPServer) HandleVersion(w http.ResponseWriter, r *http.Request) {
	b, _ := json.MarshalIndent(h.version, "", "  ")
	_, _ = w.Write(b)
}

// SetReady ...
func (h *HTTPServer) SetReady(ready bool) {
	h.ready = ready
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
