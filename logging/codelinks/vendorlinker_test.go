package codelinks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestLinker(t *testing.T) *vendorLinker {
	projectRoot := "/root/project"
	repoURL := "https://github.com/org/project3"
	vendorLinker, err := newVendorLinker(projectRoot, repoURL, "v2.0.0")
	require.NoError(t, err)
	require.NotNil(t, vendorLinker)
	assert.Equal(t, projectRoot, vendorLinker.ProjectRoot)
	return vendorLinker
}

func TestVendorLinkerProjectLink(t *testing.T) {
	linker := createTestLinker(t)
	file := "/root/project/pkg/foo.go"
	line := 42
	expected := "https://github.com/org/project3/tree/v2.0.0/pkg/foo.go#L42"
	assert.Equal(t, expected, linker.Link(file, line))
}

func TestVendorLinkerDependencyLink(t *testing.T) {
	linker := createTestLinker(t)
	file := "/root/project/vendor/github.com/org/project1/pkg/foo.go"
	line := 25
	expected := "https://github.com/org/project3/tree/v2.0.0/vendor/github.com/org/project1/pkg/foo.go#L25"
	assert.Equal(t, expected, linker.Link(file, line))
}

func TestVendorLinkerErrorUnsupportedRemote(t *testing.T) {
	linker, err := newVendorLinker("/root/project", "https://google.golang.org/project", "v1.0.0")
	assert.Error(t, err)
	assert.Nil(t, linker)
}
