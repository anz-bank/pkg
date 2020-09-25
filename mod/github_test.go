package mod

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"bou.ke/monkey"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGitHubMgrInit(t *testing.T) {
	dir := ".pkgcache"

	_, err := newGitHubMgr(GitHubOptions{CacheDir: dir})
	assert.NoError(t, err)

	_, err = newGitHubMgr(GitHubOptions{CacheDir: dir, AccessToken: accessTokenForTest(t)})
	assert.NoError(t, err)

	_, err = newGitHubMgr(GitHubOptions{})
	assert.Error(t, err)
	_, err = newGitHubMgr(GitHubOptions{AccessToken: accessTokenForTest(t)})
	assert.Error(t, err)
}

func TestGitHubMgrFind(t *testing.T) {
	cacheDir := ".pkgcache"
	repo := "github.com/foo/bar"
	tagRef, masterRef := "v0.2.0", "v0.0.0-41f04d3bba15"
	tagRepoDir := strings.Join([]string{cacheDir, repo, tagRef}, "@")
	masterRepoDir := strings.Join([]string{cacheDir, repo, masterRef}, "@")
	filea, fileb := "filea", "fileb"

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

	monkey.Patch(FileExists, func(_ afero.Fs, filename string, _ bool) bool {
		files := []string{
			filepath.Join(tagRepoDir, filea),
			filepath.Join(tagRepoDir, fileb),
			filepath.Join(masterRepoDir, filea),
		}
		for _, f := range files {
			if filename == f {
				return true
			}
		}
		return false
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(githubmod), "GetCacheRef",
		func(_ *githubMgr, _ *githubRepoPath, ref string) (string, error) {
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

func TestGitHubMgrLoad(t *testing.T) {
	cacheDir := ".pkgcache"
	githubmod := &githubMgr{cacheDir: cacheDir, fs: afero.NewMemMapFs()}

	repo := "github.com/foo/bar"
	tagRef, masterRef := "v0.2.0", "v0.0.0-41f04d3bba15"
	tagRepoDir := strings.Join([]string{repo, tagRef}, "@")
	masterRepoDir := strings.Join([]string{repo, masterRef}, "@")

	err := writeFile(githubmod.fs, filepath.Join(cacheDir, tagRepoDir, "specfile"), []byte{})
	assert.NoError(t, err)
	err = writeFile(githubmod.fs, filepath.Join(cacheDir, masterRepoDir, "specfile"), []byte{})
	assert.NoError(t, err)

	var testmods Modules
	err = githubmod.Load(&testmods)
	assert.NoError(t, err)
	assert.Equal(t, 2, testmods.Len())
	assert.Equal(t, masterRef, testmods[0].Version)
	assert.Equal(t, tagRef, testmods[1].Version)
}

func TestGitHubMgrLoadNoModules(t *testing.T) {
	cacheDir := ".pkgcache"
	fs := afero.NewMemMapFs()
	githubmod := &githubMgr{cacheDir: cacheDir, fs: fs}

	var testmods Modules
	err := githubmod.Load(&testmods)
	assert.NoError(t, err)
	assert.Equal(t, 0, testmods.Len())

	assert.True(t, (FileExists(fs, filepath.Join(cacheDir, "github.com"), true)))
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

func TestWriteFile(t *testing.T) {
	cacheDir := ".pkgcache"
	repo := "github.com/foo/bar"
	tagRef := "v0.2.0"
	tagRepoDir := strings.Join([]string{repo, tagRef}, "@")
	fs := afero.NewMemMapFs()
	content := []byte("Hello Spec!")

	err := writeFile(fs, filepath.Join(cacheDir, tagRepoDir, "specfile"), content)
	assert.NoError(t, err)
	b, err := afero.ReadFile(fs, filepath.Join(cacheDir, tagRepoDir, "specfile"))
	assert.NoError(t, err)
	assert.Equal(t, content, b)
}

func accessTokenForTest(t testing.TB) string {
	const tokenName = "GITHUB_ACCESS_TOKEN"
	token := os.Getenv(tokenName)
	if token == "" {
		t.Logf("%s empty", tokenName)
	} else {
		t.Logf("%s not empty", tokenName)
	}
	return token
}
