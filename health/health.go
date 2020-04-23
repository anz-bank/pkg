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
	pb.UnimplementedHealthServer

	ready   bool
	version *pb.VersionResponse
}

var semverRe = regexp.MustCompile(`^` +
	`v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`)

// NewServer ...
func NewServer() (*Server, error) {
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
	return &Server{version: version}, nil
}

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

// SetReady ...
func (s *Server) SetReady(ready bool) {
	s.ready = ready
}

// RegisterWith ...
func (s *Server) RegisterWith(grpcServer *grpc.Server) {
	pb.RegisterHealthServer(grpcServer, s)
}

// Alive ...
func (s *Server) Alive(_ context.Context, _ *pb.AliveRequest) (*pb.AliveResponse, error) {
	return &pb.AliveResponse{}, nil
}

// Ready ...
func (s *Server) Ready(_ context.Context, _ *pb.ReadyRequest) (*pb.ReadyResponse, error) {
	return &pb.ReadyResponse{Ready: s.ready}, nil
}

// Version ...
func (s *Server) Version(_ context.Context, _ *pb.VersionRequest) (*pb.VersionResponse, error) {
	return s.version, nil
}

// ServeHTTP ...
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		msg := fmt.Sprintf("%d method not allowed, use GET", http.StatusMethodNotAllowed)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}
	switch r.URL.Path {
	case "/healthz":
		s.HandleAlive(w, r)
	case "/readyz":
		s.HandleReady(w, r)
	case "/version":
		s.HandleVersion(w, r)
	default:
		http.NotFound(w, r)
	}
}

// HandleAlive ...
func (s *Server) HandleAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d ok\n", http.StatusOK)
}

// HandleReady ...
func (s *Server) HandleReady(w http.ResponseWriter, r *http.Request) {
	if !s.ready {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "%d service unavailable\n", http.StatusServiceUnavailable)
		return
	}
	fmt.Fprintf(w, "%d ok\n", http.StatusOK)
}

// HandleVersion ...
func (s *Server) HandleVersion(w http.ResponseWriter, r *http.Request) {
	b, _ := json.MarshalIndent(s.version, "", "  ")
	_, _ = w.Write(b)
}
