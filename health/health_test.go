//nolint: bodyclose
package health

import (
	context "context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/anz-bank/pkg/health/pb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestNewServer(t *testing.T) {
	s, err := NewServer()
	require.NoError(t, err)
	require.NotNil(t, s)
}

func TestSetReady(t *testing.T) {
	s, _ := NewServer()
	require.False(t, s.healthData.ready)
	s.SetReady(true)
	require.True(t, s.healthData.ready)
	s.SetReady(false)
	require.False(t, s.healthData.ready)
}

func TestAlive(t *testing.T) {
	s, err := NewGRPCServer()
	require.NoError(t, err)
	ctx := context.Background()
	req := &pb.AliveRequest{}
	resp, err := s.Alive(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestReady(t *testing.T) {
	s, err := NewGRPCServer()
	require.NoError(t, err)
	ctx := context.Background()
	req := &pb.ReadyRequest{}
	resp, err := s.Ready(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.False(t, resp.Ready)

	s.SetReady(true)
	resp, err = s.Ready(ctx, req)
	require.NoError(t, err)
	require.True(t, resp.Ready)
}

func resetGlobals() {
	RepoURL = Undefined
	CommitHash = Undefined
	BuildLogURL = Undefined
	ContainerTag = Undefined
	Semver = Undefined
	ScannerURLs = ""
}

func TestNewServerBadURLErr(t *testing.T) {
	defer resetGlobals()
	RepoURL = "bad URL"
	_, err := NewServer()
	require.Error(t, err)
	target := &url.Error{}
	require.True(t, errors.As(err, &target))
}

func TestNewHTTPServerErr(t *testing.T) {
	defer resetGlobals()
	RepoURL = "not a URL"
	_, err := NewHTTPServer()
	require.Error(t, err)
	target := &url.Error{}
	require.True(t, errors.As(err, &target))
}

func TestNewGRPCServerErr(t *testing.T) {
	defer resetGlobals()
	RepoURL = "not a URL"
	_, err := NewGRPCServer()
	require.Error(t, err)
	target := &url.Error{}
	require.True(t, errors.As(err, &target))
}

func TestNewServerBadJSONErr(t *testing.T) {
	defer resetGlobals()
	ScannerURLs = `{"truncated...`
	_, err := NewServer()
	require.Error(t, err)
	target := &json.SyntaxError{}
	require.True(t, errors.As(err, &target))
}

func TestVersion(t *testing.T) {
	defer resetGlobals()
	RepoURL = "http://github.com/anz-bank/pkg"
	CommitHash = "1ee4e1f233caea38d6e331299f57dd86efb47361"
	BuildLogURL = "https://github.com/anz-bank/pkg/actions/runs/84341844"
	ContainerTag = "gcr.io/google-containers/hugo"
	Semver = "v0.0.0"
	ScannerURLs = `{"scanner1" : "http://example.com", "scanner2" : "https://scan2"}`

	s, err := NewGRPCServer()
	require.NoError(t, err)

	ctx := context.Background()
	req := &pb.VersionRequest{}
	resp, err := s.Version(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, RepoURL, resp.RepoUrl)
	require.Equal(t, CommitHash, resp.CommitHash)
	require.Equal(t, BuildLogURL, resp.BuildLogUrl)
	require.Equal(t, ContainerTag, resp.ContainerTag)
	require.Equal(t, Semver, resp.Semver)
	want := map[string]string{"scanner1": "http://example.com", "scanner2": "https://scan2"}
	require.Equal(t, want, resp.ScannerUrls)
}

func versionFixture() *pb.VersionResponse {
	return &pb.VersionResponse{
		RepoUrl:      "http://github.com/anz-bank/pkg",
		CommitHash:   "1ee4e1f233caea38d6e331299f57dd86efb47361",
		BuildLogUrl:  "https://github.com/anz-bank/pkg/actions/runs/84341844",
		ContainerTag: "gcr.io/google-containers/hugo",
		Semver:       "v0.0.0",
		ScannerUrls: map[string]string{
			"scanner1": "http://example.com",
			"scanner2": "https://scan2",
		},
	}
}

func TestValidationBadScannerURL(t *testing.T) {
	v := versionFixture()
	v.ScannerUrls["badScanner"] = "bad scanner url"
	err := validateVersion(v)
	require.Error(t, err)
	target := &url.Error{}
	require.True(t, errors.As(err, &target))
}

func TestValidationBadBuildLogURL(t *testing.T) {
	v := versionFixture()
	v.BuildLogUrl = "bad buildlog url"
	err := validateVersion(v)
	require.Error(t, err)
	target := &url.Error{}
	require.True(t, errors.As(err, &target))
}

func TestValidationBadSemver(t *testing.T) {
	v := versionFixture()
	v.Semver = "bad semver"
	err := validateVersion(v)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrInvalidSemver))
}

func TestHTTPAlive(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/healthz", nil)
	w := httptest.NewRecorder()
	s, err := NewHTTPServer()
	require.NoError(t, err)

	s.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "200 ok\n", string(body))
}

func TestHTTPNotReady(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/readyz", nil)
	w := httptest.NewRecorder()
	s, err := NewHTTPServer()
	require.NoError(t, err)

	s.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	require.Equal(t, "503 service unavailable\n", string(body))
}

func TestHTTPReady(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/readyz", nil)
	w := httptest.NewRecorder()
	s, err := NewHTTPServer()
	require.NoError(t, err)
	s.SetReady(true)

	s.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "200 ok\n", string(body))
}

func TestHTTPVersion(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/version", nil)
	w := httptest.NewRecorder()
	s, err := NewHTTPServer()
	require.NoError(t, err)
	s.healthData.version = versionFixture()

	s.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	want, _ := json.Marshal(versionFixture())
	require.JSONEq(t, string(want), string(body))
}

func TestHTTPNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/MISSING_PATH", nil)
	w := httptest.NewRecorder()
	s, err := NewHTTPServer()
	require.NoError(t, err)
	s.healthData.version = versionFixture()

	s.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
	require.Equal(t, "404 page not found\n", string(body))
}

func TestHTTPMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest("POST", "http://example.com/MISSING_PATH", nil)
	w := httptest.NewRecorder()
	s, err := NewHTTPServer()
	require.NoError(t, err)
	s.healthData.version = versionFixture()

	s.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	require.Equal(t, "405 method not allowed, use GET\n", string(body))
}

func TestRegisterWith(t *testing.T) {
	grpcServer := grpc.NewServer()
	hs, err := NewGRPCServer()
	require.NoError(t, err)
	hs.RegisterWith(grpcServer)
	info := grpcServer.GetServiceInfo()
	_, ok := info["anz.health.v1.Health"]
	require.True(t, ok)
}
