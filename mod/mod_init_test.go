package mod

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSyslModInit(t *testing.T) {
	fs := afero.NewOsFs()

	// assumes the test folder (cwd) is not a go module folder
	removeFile(t, fs, "go.sum")
	removeFile(t, fs, "go.mod")

	err := SyslModInit("test")
	assert.NoError(t, err)

	removeFile(t, fs, "go.sum")
	removeFile(t, fs, "go.mod")
}

func TestSyslModInitAlreadyExists(t *testing.T) {
	fs := afero.NewOsFs()

	// assumes the test folder (cwd) is not a go module folder
	removeFile(t, fs, "go.sum")
	removeFile(t, fs, "go.mod")

	err := SyslModInit("test")
	assert.NoError(t, err)

	err = SyslModInit("test")
	assert.Error(t, err)

	removeFile(t, fs, "go.sum")
	removeFile(t, fs, "go.mod")
}
