package mod

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestModInit(t *testing.T) {
	fs := afero.NewOsFs()

	// assumes the test folder (cwd) is not a go module folder
	removeFile(t, fs, "go.sum")
	removeFile(t, fs, "go.mod")

	err := ModInit("test")
	assert.NoError(t, err)

	removeFile(t, fs, "go.sum")
	removeFile(t, fs, "go.mod")
}

func TestModInitAlreadyExists(t *testing.T) {
	fs := afero.NewOsFs()

	// assumes the test folder (cwd) is not a go module folder
	removeFile(t, fs, "go.sum")
	removeFile(t, fs, "go.mod")

	err := ModInit("test")
	assert.NoError(t, err)

	err = ModInit("test")
	assert.Error(t, err)

	removeFile(t, fs, "go.sum")
	removeFile(t, fs, "go.mod")
}
