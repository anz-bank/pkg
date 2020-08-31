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

func TestAppendVersion(t *testing.T) {
	assert.Equal(t, "github.com/anz-bank/pkg@v0.1", AppendVersion("github.com/anz-bank/pkg", "v0.1"))
	assert.Equal(t, "github.com/anz-bank/pkg/foo@v0.2", AppendVersion("github.com/anz-bank/pkg/foo", "v0.2"))
	assert.Equal(t, "github.com/anz-bank/pkg/foo", AppendVersion("github.com/anz-bank/pkg/foo", ""))
}
