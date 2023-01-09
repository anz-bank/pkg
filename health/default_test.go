package health

import (
	"errors"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestDefaultHTTPSetReady(t *testing.T) {
	resetDefaults()
	defer resetDefaults()

	require.False(t, defaultState.IsReady())
	require.Nil(t, DefaultServer)

	SetReady(true)
	require.True(t, defaultState.IsReady())
	require.Nil(t, DefaultServer)

	mux := http.NewServeMux()

	err := RegisterWithHTTP(mux)
	require.NoError(t, err)
	require.True(t, defaultState.IsReady())
	require.NotNil(t, DefaultServer)
	require.True(t, DefaultServer.State.IsReady())
}

func TestDefaultGRPCSetReady(t *testing.T) {
	resetDefaults()
	defer resetDefaults()

	require.False(t, defaultState.IsReady())
	require.Nil(t, DefaultServer)

	SetReady(true)
	require.True(t, defaultState.IsReady())
	require.Nil(t, DefaultServer)

	grpcServer := grpc.NewServer()

	err := RegisterWithGRPC(grpcServer)
	require.NoError(t, err)
	require.True(t, defaultState.IsReady())
	require.NotNil(t, DefaultServer)
	require.True(t, DefaultServer.State.IsReady())
}

func TestDefaultHTTPSetReadyErr(t *testing.T) {
	resetDefaults()
	defer resetDefaults()

	RepoURL = "not a URL"
	mux := http.NewServeMux()
	err := RegisterWithHTTP(mux)
	require.Error(t, err)
	target := &url.Error{}
	require.True(t, errors.As(err, &target))
}

func TestDefaultGRPCSetReadyErr(t *testing.T) {
	resetDefaults()
	defer resetDefaults()

	RepoURL = "not a URL"
	grpcServer := grpc.NewServer()
	err := RegisterWithGRPC(grpcServer)
	require.Error(t, err)
	target := &url.Error{}
	require.True(t, errors.As(err, &target))
}

func TestSetReadyProvider(t *testing.T) {
	resetDefaults()
	defer resetDefaults()

	var r readiness
	r.SetReady(true)
	SetReadyProvider(&r)
	require.Nil(t, DefaultServer)
	mux := http.NewServeMux()
	err := RegisterWithHTTP(mux)
	require.NoError(t, err)
	require.True(t, DefaultServer.State.IsReady())
}

func TestRaceConditionWhenRegisteringServices(t *testing.T) {
	resetDefaults()
	defer resetDefaults()
	var r readiness
	r.SetReady(true)
	SetReadyProvider(&r)
	require.Nil(t, DefaultServer)

	mux := http.NewServeMux()
	err := RegisterWithHTTP(mux)
	require.NoError(t, err)

	DefaultServer = nil // simulate the first call to newDefaultServer taking too long or failing entirely

	grpcServer := grpc.NewServer()
	err = RegisterWithGRPC(grpcServer)
	require.Error(t, err)
}

func resetDefaults() {
	DefaultServer = nil
	defaultState = State{ReadyProvider: new(readiness)}
	serverInit = sync.Once{}

	resetGlobals()
}
