package mod

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestGitHubMgrInit(t *testing.T) {
	githubmod := &githubMgr{}
	dir := ".pkgcache"
	err := githubmod.Init(&dir, nil)
	assert.NoError(t, err)

	var accessToken *string
	rawToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	if rawToken != "" {
		accessToken = &rawToken
	}
	err = githubmod.Init(&dir, accessToken)
	assert.NoError(t, err)

	err = githubmod.Init(nil, nil)
	assert.Error(t, err)
	err = githubmod.Init(nil, accessToken)
	assert.Error(t, err)
}

func TestGitHubMgrGet(t *testing.T) {
	githubmod := &githubMgr{}
	dir := ".pkgcache"
	var accessToken *string
	rawToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	if rawToken != "" {
		accessToken = &rawToken
	}
	err := githubmod.Init(&dir, accessToken)
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

func TestGitHubMgrFind(t *testing.T) {
	cacheDir := ".pkgcache"
	repo := "github.com/foo/bar"
	filea := "filea"
	fileb := "fileb"
	tagRef := "v0.2.0"
	masterRef := "v0.0.0-41f04d3bba15"
	tagRepoDir := "github.com/foo/bar@v0.2.0"
	masterRepoDir := "github.com/foo/bar@v0.0.0-41f04d3bba15"

	githubmod := &githubMgr{cacheDir: cacheDir}
	testMods := Modules{}
	tagMod := &Module{
		Name:    repo,
		Version: tagRef,
		Dir:     tagRepoDir,
	}
	masterMod := &Module{
		Name:    repo,
		Version: masterRef,
		Dir:     masterRepoDir,
	}
	testMods.Add(tagMod)
	testMods.Add(masterMod)

	monkey.Patch(FileExists, func(filename string, _ bool) bool {
		files := []string{
			filepath.Join(cacheDir, tagRepoDir, filea),
			filepath.Join(cacheDir, tagRepoDir, fileb),
			filepath.Join(cacheDir, masterRepoDir, filea),
		}
		for _, f := range files {
			if filename == f {
				return true
			}
		}
		return false
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(githubmod), "CacheRef", func(_ *githubMgr, _ *githubRepoPath, ref string) (string, error) {
		switch ref {
		case tagRef:
			return tagRef, nil
		case MasterBranch:
			return masterRef, nil
		}
		return "", fmt.Errorf("ref not found")
	})
	defer monkey.UnpatchAll()

	assert.Equal(t, tagMod, githubmod.Find(path.Join(repo, filea), tagRef, &testMods))
	assert.Equal(t, tagMod, githubmod.Find(path.Join(repo, fileb), tagRef, &testMods))
	assert.Nil(t, githubmod.Find(repo, tagRef, &testMods))
	assert.Nil(t, githubmod.Find(path.Join(repo, "wrong"), tagRef, &testMods))

	assert.Equal(t, masterMod, githubmod.Find(path.Join(repo, filea), MasterBranch, &testMods))
	assert.Equal(t, masterMod, githubmod.Find(path.Join(repo, filea), "", &testMods))
	assert.Nil(t, githubmod.Find(repo, MasterBranch, &testMods))
	assert.Nil(t, githubmod.Find(path.Join(repo, fileb), MasterBranch, &testMods))

	assert.Nil(t, githubmod.Find("github.com/foo/wrongrepo/files", tagRef, &testMods))
}

func TestGetGitHubRepoPath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		filename string
		path     *githubRepoPath
	}{
		{"github.com/anz-bank/pkg", nil},
		{"github.com/anz-bank/pkg/", nil},
		{"github.com/anz-bank/pkg/deps.sysl", &githubRepoPath{"anz-bank", "pkg", "deps.sysl"}},
		{"github.com/anz-bank/pkg/nested/module/deps.sysl", &githubRepoPath{"anz-bank", "pkg", "nested/module/deps.sysl"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.filename, func(t *testing.T) {
			t.Parallel()
			p, err := getGitHubRepoPath(tt.filename)
			if tt.path == nil {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.path, p)
		})
	}
}

func TestCacheRef(t *testing.T) {
	githubmod := &githubMgr{}
	dir := ".pkgcache"
	err := githubmod.Init(&dir, nil)
	assert.NoError(t, err)
	repoPath := &githubRepoPath{
		owner: "anz-bank",
		repo:  "pkg",
	}
	ref, err := githubmod.CacheRef(repoPath, "v0.0.7")
	assert.NoError(t, err)
	assert.Equal(t, "v0.0.7", ref)

	ref, err = githubmod.CacheRef(repoPath, MasterBranch)
	assert.NoError(t, err)
	assert.Equal(t, "v0.0.0-", ref[:7])
}
