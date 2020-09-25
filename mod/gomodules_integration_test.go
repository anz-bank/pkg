//+build integration

package mod

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoModulesGet(t *testing.T) {
	gomod := &goModules{}
	testMods := Modules{}

	mod, err := gomod.Get(RemoteDepsFile, "", &testMods)
	assert.Nil(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)

	mod, err = gomod.Get(RemoteDepsFile, MasterBranch, &testMods)
	assert.Nil(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)

	mod, err = gomod.Get(RemoteDepsFile, "v0.0.1", &testMods)
	assert.Nil(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)
	assert.Equal(t, "v0.0.1", mod.Version)

	mod, err = gomod.Get("github.com/anz-bank/wrongpath", "", &testMods)
	assert.Error(t, err)
	assert.Nil(t, mod)
}
