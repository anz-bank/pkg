package mod

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/go-github/v32/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"golang.org/x/oauth2"
)

type githubMgr struct {
	client   *github.Client
	cacheDir string
	fs       afero.Fs
}

type GitHubOptions struct {
	CacheDir    string
	AccessToken string
	Fs          afero.Fs
}

func newGitHubMgr(opt GitHubOptions) (*githubMgr, error) {
	d := &githubMgr{}
	if opt.AccessToken == "" {
		d.client = github.NewClient(nil)
	} else {
		// Authenticated clients can make up to 5,000 requests per hour.
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: opt.AccessToken},
		)
		tc := oauth2.NewClient(context.Background(), ts)

		d.client = github.NewClient(tc)
	}

	if opt.CacheDir == "" {
		return nil, errors.New("cache directory cannot be empty")
	}
	d.cacheDir = opt.CacheDir

	if opt.Fs != nil {
		d.fs = opt.Fs
	} else {
		d.fs = afero.NewOsFs()
	}
	return d, nil
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

type RateLimitError = github.RateLimitError

func (d *githubMgr) Get(filename, ver string, m *Modules) (*Module, error) {
	repoPath, err := getGitHubRepoPath(filename)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	if ver == "" {
		ver = MasterBranch
	}
	refOps := &github.RepositoryContentGetOptions{Ref: ver}

	fileContent, _, _, err := d.client.Repositories.GetContents(ctx, repoPath.owner, repoPath.repo, repoPath.path, refOps)
	if err != nil {
		if err, ok := err.(*github.RateLimitError); ok {
			return nil, err
		}
		if err, ok := err.(*github.ErrorResponse); ok && err.Response.StatusCode == http.StatusNotFound {
			return nil, &NotFoundError{Message: err.Error()}
		}
		return nil, err
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return nil, err
	}

	ref, err := d.GetCacheRef(repoPath, ver)
	if err != nil {
		return nil, err
	}

	name := strings.Join([]string{"github.com", repoPath.owner, repoPath.repo}, "/")
	dir := filepath.Join(d.cacheDir, "github.com", repoPath.owner, repoPath.repo)
	dir = AppendVersion(dir, ref)
	new := &Module{
		Name:    name,
		Dir:     dir,
		Version: ref,
	}

	fname := filepath.Join(dir, repoPath.path)
	if !FileExists(d.fs, fname, false) {
		err = writeFile(d.fs, fname, []byte(content))
		if err != nil {
			return nil, err
		}
		m.Add(new)
	}

	return new, nil
}

func (d *githubMgr) Find(filename, ver string, m *Modules) *Module {
	if ver == "" {
		ver = MasterBranch
	}

	repoPath, err := getGitHubRepoPath(filename)
	if err != nil {
		logrus.Debug("get github repository path error:", err)
		return nil
	}

	ref, err := d.GetCacheRef(repoPath, ver)
	if err != nil {
		logrus.Debug("get github repository ref error:", err)
		return nil
	}

	for _, mod := range *m {
		if hasPathPrefix(mod.Name, filename) {
			if mod.Version == ref {
				relpath, err := filepath.Rel(mod.Name, filename)
				if err == nil && FileExists(d.fs, filepath.Join(mod.Dir, relpath), false) {
					return mod
				}
			}
		}
	}

	return nil
}

func (d *githubMgr) Load(m *Modules) error {
	githubPath := filepath.Join(d.cacheDir, "github.com")
	if !FileExists(d.fs, githubPath, true) {
		if err := d.fs.MkdirAll(githubPath, 0770); err != nil {
			return err
		}
		return nil
	}

	githubDir, err := d.fs.Open(githubPath)
	if err != nil {
		return err
	}

	owners, err := githubDir.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, owner := range owners {
		ownerDir, err := d.fs.Open(filepath.Join(githubPath, owner))
		if err != nil {
			return err
		}
		repos, err := ownerDir.Readdirnames(-1)
		if err != nil {
			return err
		}
		for _, repo := range repos {
			p, ver := ExtractVersion(repo)
			name := filepath.Join("github.com", owner, p)
			m.Add(&Module{
				Name:    name,
				Dir:     filepath.Join(ownerDir.Name(), repo),
				Version: ver,
			})
		}
	}

	return nil
}

type githubRepoPath struct {
	owner string
	repo  string
	path  string
}

func getGitHubRepoPath(filename string) (*githubRepoPath, error) {
	names := strings.FieldsFunc(filename, func(c rune) bool {
		return c == '/'
	})
	if len(names) < 4 {
		return nil, fmt.Errorf("the imported module path %s is invalid", filename)
	}
	if names[0] != "github.com" {
		return nil, errors.New("non-github.com repository is not supported under GitHub mode")
	}

	owner := names[1]
	repo := names[2]
	path := path.Join(names[3:]...)

	return &githubRepoPath{
		owner: owner,
		repo:  repo,
		path:  path,
	}, nil
}

func writeFile(fs afero.Fs, filename string, content []byte) error {
	if err := fs.MkdirAll(filepath.Dir(filename), 0770); err != nil {
		return err
	}
	file, err := fs.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Write(content); err != nil {
		return err
	}
	return nil
}

const SHALength = 12

func (d *githubMgr) GetCacheRef(repoPath *githubRepoPath, ref string) (string, error) {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	resolvedRef := make(chan string, 2)
	errors := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		branch, _, err := d.client.Git.GetRef(ctx, repoPath.owner, repoPath.repo, "heads/"+ref)
		if err != nil {
			errors <- err
		} else {
			resolvedRef <- "v0.0.0-" + branch.GetObject().GetSHA()[:SHALength]
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _, err := d.client.Git.GetRef(ctx, repoPath.owner, repoPath.repo, "tags/"+ref)
		if err != nil {
			errors <- err
		} else {
			resolvedRef <- ref
		}
	}()

	wg.Wait()
	if len(errors) >= 2 {
		return "", fmt.Errorf("failed to find cache ref %w", <-errors)
	} else if len(resolvedRef) == 2 {
		return "", fmt.Errorf("ref is both tag and branch")
	}
	return <-resolvedRef, nil
}
