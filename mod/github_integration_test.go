//+build integration

package mod

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGitHubMgrGet(t *testing.T) {
	t.Parallel()
	dir := ".pkgcache"
	fs := afero.NewMemMapFs()
	githubmod, err := newGitHubMgr(GitHubOptions{CacheDir: dir, AccessToken: accessTokenForTest(t), Fs: fs})
	assert.NoError(t, err)
	testMods := Modules{}

	mod, err := githubmod.Get(RemoteDepsFile, "", &testMods)
	assert.Nil(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)

	mod, err = githubmod.Get(RemoteDepsFile, MasterBranch, &testMods)
	assert.Nil(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)

	mod, err = githubmod.Get(RemoteDepsFile, "v0.0.1", &testMods)
	assert.Nil(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)
	assert.Equal(t, "v0.0.1", mod.Version)

	mod, err = githubmod.Get("github.com/anz-bank/wrong/path", "", &testMods)
	assert.Error(t, err)
	assert.Nil(t, mod)
}

func TestGetCacheRef(t *testing.T) {
	t.Parallel()
	dir := ".pkgcache"
	githubmod, err := newGitHubMgr(GitHubOptions{CacheDir: dir, AccessToken: accessTokenForTest(t)})
	assert.NoError(t, err)
	repoPath := &githubRepoPath{
		owner: "anz-bank",
		repo:  "pkg",
	}
	ref, err := githubmod.GetCacheRef(repoPath, "v0.0.7")
	assert.NoError(t, err)
	assert.Equal(t, "v0.0.7", ref)

	ref, err = githubmod.GetCacheRef(repoPath, MasterBranch)
	assert.NoError(t, err)
	assert.Equal(t, "v0.0.0-", ref[:7])
}

func BenchmarkGetCacheRef(b *testing.B) {
	dir := ".pkgcache"
	githubmod, _ := newGitHubMgr(GitHubOptions{CacheDir: dir, AccessToken: accessTokenForTest(b)})
	repoPath := &githubRepoPath{
		owner: "anz-bank",
		repo:  "pkg",
	}
	for i := 0; i < b.N; i++ {
		_, _ = githubmod.GetCacheRef(repoPath, "v0.0.7")
	}
}
