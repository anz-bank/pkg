// +build integration

package mod

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	SyslDepsFile   = "github.com/anz-bank/sysl/tests/deps.sysl"
	SyslRepo       = "github.com/anz-bank/sysl"
	RemoteDepsFile = "github.com/anz-bank/sysl-examples/demos/simple/simple.sysl"
	RemoteRepo     = "github.com/anz-bank/sysl-examples"
)

func TestRetrieveGoModules(t *testing.T) {
	fs := afero.NewOsFs()
	createGomodFile(t, fs)
	defer removeGomodFile(t, fs)

	filename := SyslDepsFile
	mod, err := Retrieve(filename, "")
	require.NoError(t, err)
	assert.Equal(t, SyslRepo, mod.Name)

	filename = RemoteDepsFile
	mod, err = Retrieve(filename, "")
	require.NoError(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)

	mod, err = Retrieve(filename, "v0.0.1")
	require.NoError(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)
	assert.Equal(t, "v0.0.1", mod.Version)
}

func TestRetrieveGitHubMode(t *testing.T) {
	mode.modeType = GitHubMode
	defer func() {
		mode.modeType = GoModulesMode
	}()

	filename := SyslDepsFile
	mod, err := Retrieve(filename, "")
	require.NoError(t, err)
	assert.Equal(t, SyslRepo, mod.Name)

	filename = RemoteDepsFile
	mod, err = Retrieve(filename, "")
	require.NoError(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)

	mod, err = Retrieve(filename, "v0.0.1")
	require.NoError(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)
	assert.Equal(t, "v0.0.1", mod.Version)
}

func BenchmarkRetrieveGitHubModeCached(b *testing.B) {
	mode.modeType = GitHubMode
	defer func() {
		mode.modeType = GoModulesMode
	}()
	dir := ".pkgcache"
	_ = Config(GitHubMode, GoModulesOptions{},
		GitHubOptions{CacheDir: dir, AccessToken: accessTokenForTest(b), Fs: afero.NewMemMapFs()})

	// Fetch files once to cache
	_, _ = Retrieve(SyslDepsFile, "")
	_, _ = Retrieve(RemoteDepsFile, "")
	_, _ = Retrieve(RemoteDepsFile, "v0.0.1")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Retrieve(SyslDepsFile, "")
		_, _ = Retrieve(RemoteDepsFile, "")
		_, _ = Retrieve(RemoteDepsFile, "v0.0.1")
	}
}
