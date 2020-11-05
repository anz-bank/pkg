package codelinks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfoFromPath(t *testing.T) {
	modPath := "/gopath/pkg/mod"
	fullPath := "/gopath/pkg/mod/host/path/to/repo@v1.0.0/path/to/file.go"
	module, file, version := infoFromPath(modPath, fullPath)

	assert.Equal(t, "host/path/to/repo", module)
	assert.Equal(t, "path/to/file.go", file)
	assert.Equal(t, "v1.0.0", version)
}

func TestModLinkerErrorForUnsupportedRemote(t *testing.T) {
	linker, err := newModLinker("/root/my/project", "/root/go/pkg/mod", "https://google.golang.org/project", "v1.0.0")
	assert.Error(t, err)
	assert.Nil(t, linker)
}

func TestModLinkerProjectLink(t *testing.T) {
	linker, err := newModLinker("/root/my/project", "/root/go/pkg/mod", "https://github.com/myorg/myrepo", "v1.0.0")
	require.NoError(t, err)
	require.NotNil(t, linker)
	file := "/root/my/project/path/to/file.go"
	line := 25

	expected := "https://github.com/myorg/myrepo/tree/v1.0.0/path/to/file.go#L25"
	actual := linker.Link(file, line)
	assert.Equal(t, expected, actual)
}

func TestModLinkerModLink(t *testing.T) {
	linker, err := newModLinker("/root/my/project", "/root/go/pkg/mod", "https://github.com/myorg/myrepo", "v1.0.0")
	require.NoError(t, err)
	require.NotNil(t, linker)
	file := "/root/go/pkg/mod/github.com/repo/path@v1.0.0/path/to/file.go"
	line := 42

	expected := "https://github.com/repo/path/tree/v1.0.0/path/to/file.go#L42"
	actual := linker.Link(file, line)
	assert.Equal(t, expected, actual)
}

func TestModLinkerModuleNotInModpath(t *testing.T) {
	linker, err := newModLinker("/root/my/project", "/root/go/pkg/mod", "https://github.com/myorg/myrepo", "v1.0.0")
	require.NoError(t, err)
	require.NotNil(t, linker)
	file := "/root/notgopath/pkg/mod/github.com/repo/path@v1.0.0/path/to/file.go"
	line := 42

	actual := linker.Link(file, line)
	assert.Equal(t, notFound, actual)
}

func TestModuleInfoExtractsCommitSha(t *testing.T) {
	modpath := "/root/go/pkg/mod"
	file := "/root/go/pkg/mod/host/org/repo@v0.0.0-crap-crap-12345/file.go"
	_, _, version := infoFromPath(modpath, file)
	assert.Equal(t, "12345", version)
}
