package codelinks

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testRepo    = "https://github.com/anz-bank/pkg/logging"
	testVersion = "v0.1.0"
)

func TestGetCodeLinker(t *testing.T) {
	linker, err := GetCodeLinker()
	require.NoError(t, err)
	require.NotNil(t, linker)
	assert.IsType(t, LocalLinker{}, linker)
}

func TestGetCodeLinkerWithRemotes(t *testing.T) {
	projectRoot, err := os.Getwd()
	require.NoError(t, err)

	// We fudge the gopath here to trick the is-built-from-vendor check
	modpath := projectRoot

	linker, err := getCodeLinker(codeLinkerConfig{
		localLinks:  false,
		fromVendor:  false,
		projectRoot: projectRoot,
		modpath:     modpath,
		repoURL:     testRepo,
		version:     testVersion,
	})
	require.NoError(t, err)
	require.NotNil(t, linker)
	require.IsType(t, &modLinker{}, linker)

	modlinker := linker.(*modLinker)
	assert.Equal(t, modpath, modlinker.ModPath)
	assert.Equal(t, projectRoot, modlinker.ProjectRoot)
}

func TestGetCodeLinkerVendor(t *testing.T) {
	projectRoot, err := os.Getwd()
	require.NoError(t, err)

	linker, err := getCodeLinker(codeLinkerConfig{
		localLinks:  false,
		fromVendor:  true,
		projectRoot: projectRoot,
		modpath:     os.Getenv("GOPATH"),
		repoURL:     testRepo,
		version:     testVersion,
	})
	require.NoError(t, err)
	require.NotNil(t, linker)
	assert.IsType(t, &vendorLinker{}, linker)
}

func TestGetCodeLinkerErrorUnsupportedCases(t *testing.T) {
	projectRoot, err := os.Getwd()
	require.NoError(t, err)

	linker, err := getCodeLinker(codeLinkerConfig{
		localLinks:  false,
		fromVendor:  false,
		projectRoot: projectRoot,
		modpath:     "",
		repoURL:     testRepo,
		version:     testVersion,
	})

	require.Error(t, err)
	require.Nil(t, linker)
}

func TestGithubRepoFromURLNoScheme(t *testing.T) {
	url := "github.com/org/repo"
	repo := githubRepoFromURL(url)
	assert.Equal(t, url, repo)
}

func TestLocalLinker(t *testing.T) {
	file := "file"
	line := 12
	linker := LocalLinker{}
	expected := "file:12"
	actual := linker.Link(file, line)
	assert.Equal(t, expected, actual)
}

func TestCreateLinkUnsupportedRemote(t *testing.T) {
	module := "golang.org/x/sync"
	file := "pkg/foo.go"
	version := "v1.0.0"
	line := 32
	assert.Equal(t, notFound, createLink(module, file, version, line))
}

func TestGithubLinkModuleNotLengthThree(t *testing.T) {
	module := "github.com/org/repo/submodule"
	file := "pkg/foo.go"
	version := "submodule/v0.1.0"
	line := 100
	expected := "https://github.com/org/repo/tree/submodule/v0.1.0/submodule/pkg/foo.go#L100"
	actual := GithubLink(module, file, version, line)
	assert.Equal(t, expected, actual)
}

var blink string

func BenchmarkModLinker(b *testing.B) {
	linker, err := newModLinker("/root/my/project", "/root/go", "https://github.com/myorg/myrepo", "v1.0.0")
	if err != nil {
		b.Fatal(err)
	}
	file := "/root/go/pkg/mod/github.com/myorg/myrepo@v1.0.0/path/to/file.go"
	line := 42

	var link string
	for n := 0; n < b.N; n++ {
		link = linker.Link(file, line)
	}
	blink = link
}

func BenchmarkVendorLinker(b *testing.B) {
	linker, err := newVendorLinker("/root/my/project", "https://github.com/myorg/myrepo", "v1.0.0")
	if err != nil {
		b.FailNow()
	}
	file := "/root/my/project/vendor/github.com/org/repo/pkg/foo.go"
	line := 42

	var link string
	for n := 0; n < b.N; n++ {
		link = linker.Link(file, line)
	}
	blink = link
}
