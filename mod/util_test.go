package mod

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractVersion(t *testing.T) {
	path, ver := ExtractVersion("github.com/anz-bank/pkg@v0.1")
	assert.Equal(t, "github.com/anz-bank/pkg", path)
	assert.Equal(t, "v0.1", ver)

	path, ver = ExtractVersion("github.com/anz-bank/pkg/foo@v0.2")
	assert.Equal(t, "github.com/anz-bank/pkg/foo", path)
	assert.Equal(t, "v0.2", ver)

	path, ver = ExtractVersion("github.com/anz-bank/pkg/foo")
	assert.Equal(t, "github.com/anz-bank/pkg/foo", path)
	assert.Equal(t, "", ver)
}
