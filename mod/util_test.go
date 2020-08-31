package mod

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractVersion(t *testing.T) {
	path, ver := ExtractVersion("github.com/anz-bank/sysl@v0.1")
	assert.Equal(t, "github.com/anz-bank/sysl", path)
	assert.Equal(t, "v0.1", ver)

	path, ver = ExtractVersion("github.com/anz-bank/sysl/pkg@v0.2")
	assert.Equal(t, "github.com/anz-bank/sysl/pkg", path)
	assert.Equal(t, "v0.2", ver)

	path, ver = ExtractVersion("github.com/anz-bank/sysl/pkg")
	assert.Equal(t, "github.com/anz-bank/sysl/pkg", path)
	assert.Equal(t, "", ver)
}
