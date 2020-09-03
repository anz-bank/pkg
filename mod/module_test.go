package mod

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

const (
	SyslDepsFile   = "github.com/anz-bank/sysl/tests/deps.sysl"
	SyslRepo       = "github.com/anz-bank/sysl"
	RemoteDepsFile = "github.com/anz-bank/sysl-examples/demos/simple/simple.sysl"
	RemoteRepo     = "github.com/anz-bank/sysl-examples"
)

func TestConfigGitHubMode(t *testing.T) {
	dir := ".pkgcache"
	err := Config(GitHubMode, GoModulesOptions{}, GitHubOptions{CacheDir: dir, AccessToken: accessTokenForTest()})
	assert.NoError(t, err)

	err = Config(GitHubMode, GoModulesOptions{}, GitHubOptions{AccessToken: accessTokenForTest()})
	assert.Error(t, err)
}

func TestConfigGoModulesMode(t *testing.T) {
	fs := afero.NewMemMapFs()
	createGomodFile(t, fs)

	err := Config(GitHubMode, GoModulesOptions{ModName: "mod"}, GitHubOptions{CacheDir: ".pkgcache", AccessToken: accessTokenForTest()})
	assert.NoError(t, err)
}

func TestConfigWrongMode(t *testing.T) {
	err := Config("wrong", GoModulesOptions{ModName: "mod"}, GitHubOptions{AccessToken: accessTokenForTest()})
	assert.Error(t, err)
}

func TestAdd(t *testing.T) {
	var testMods Modules
	testMods.Add(&Module{Name: "modulepath"})
	assert.Equal(t, 1, len(testMods))
	assert.Equal(t, &Module{Name: "modulepath"}, testMods[0])
}

func TestLen(t *testing.T) {
	var testMods Modules
	assert.Equal(t, 0, testMods.Len())
	testMods.Add(&Module{Name: "modulepath"})
	assert.Equal(t, 1, testMods.Len())
}

func TestRetrieveGoModules(t *testing.T) {
	fs := afero.NewMemMapFs()
	createGomodFile(t, fs)

	filename := SyslDepsFile
	mod, err := Retrieve(filename, "")
	assert.NoError(t, err)
	assert.Equal(t, SyslRepo, mod.Name)

	filename = RemoteDepsFile
	mod, err = Retrieve(filename, "")
	assert.NoError(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)

	mod, err = Retrieve(filename, "v0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)
	assert.Equal(t, "v0.0.1", mod.Version)
}

func TestRetrieveWithWrongPath(t *testing.T) {
	fs := afero.NewMemMapFs()
	createGomodFile(t, fs)

	wrongpath := "wrong_file_path/deps.sysl"
	mod, err := Retrieve(wrongpath, "")
	assert.Error(t, err)
	assert.Nil(t, mod)
}

func TestRetrieveGitHubMode(t *testing.T) {
	mode = GitHubMode
	defer func() {
		mode = GoModulesMode
	}()

	filename := SyslDepsFile
	mod, err := Retrieve(filename, "")
	assert.NoError(t, err)
	assert.Equal(t, SyslRepo, mod.Name)

	filename = RemoteDepsFile
	mod, err = Retrieve(filename, "")
	assert.NoError(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)

	mod, err = Retrieve(filename, "v0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)
	assert.Equal(t, "v0.0.1", mod.Version)
}

func TestRetrieveWithWrongPathGitHubMode(t *testing.T) {
	mode = GitHubMode
	defer func() {
		mode = GoModulesMode
	}()

	wrongpath := "wrong_file_path/deps.sysl"
	mod, err := Retrieve(wrongpath, "")
	assert.Error(t, err)
	assert.Nil(t, mod)
}

func TestHasPathPrefix(t *testing.T) {
	t.Parallel()
	tests := []struct {
		prefix string
	}{
		{"github.com/anz-bank/sysl"},
		{"github.com/anz-bank/sysl/"},
		{"github.com/anz-bank/sysl/deps.sysl"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.prefix, func(t *testing.T) {
			t.Parallel()
			assert.True(t, hasPathPrefix(tt.prefix, "github.com/anz-bank/sysl/deps.sysl"))
		})
	}

	assert.False(t, hasPathPrefix("github.com/anz-bank/sysl2", "github.com/anz-bank/sysl/deps.sysl"))
}
