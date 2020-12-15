package codelinks

import (
	"errors"
	"strings"
)

var (
	ErrUnsupportedRemoteRepository = errors.New("create logging CodeLinker: Unsupported remote repository, must be github (or github enterprise)") //nolint:lll
)

type modLinker struct {
	ProjectRoot string
	ModPath     string
	projectLink func(string, int) string
}

// NewModLinker produces a ModLinker
//
// projectRoot is the root directory of the project in the build file system
//
// repoURL is the remote repository hosting the code
//
// projectVersion is the version of the project being built (use commit-sha if no tag)
//
// GoPath is the gopath of the build file system (used to find GOPATH/pkg/mod)
func newModLinker(projectRoot, modpath, repoURL, projectVersion string) (*modLinker, error) {
	var projectLink func(string, int) string
	if strings.Contains(repoURL, "github") {
		repo := githubRepoFromURL(repoURL)
		fileTrim := len(projectRoot) + 1
		projectLink = func(file string, line int) string {
			return GithubLink(repo, file[fileTrim:], projectVersion, line)
		}
	} else {
		return nil, ErrUnsupportedRemoteRepository
	}
	return &modLinker{
		ProjectRoot: projectRoot,
		ModPath:     modpath,
		projectLink: projectLink,
	}, nil
}

var notFound = "!ERR(source unknown)"

// Link produces a link for either the project or module dependency
func (l *modLinker) Link(file string, line int) string {
	if strings.HasPrefix(file, l.ProjectRoot) {
		// TODO(cantosd): handle vendored build
		return l.projectLink(file, line)
	}
	// Assume module from gopath/pkg/mod
	if !strings.HasPrefix(file, l.ModPath) {
		return notFound
	}

	module, file, version := infoFromPath(l.ModPath, file)
	return createLink(module, file, version, line)
}

func githubRepoFromURL(url string) string {
	if !strings.Contains(url, "://") {
		return url
	}
	return strings.Split(url, "://")[1]
}

func infoFromPath(modPath, fullFile string) (module string, file string, version string) {
	// Module path of the form GOPATH/pkg/mod/module/import/path@<ref>/file/path
	// Procedure is as follows
	// 1. Strip GOPATH/pkg/mod
	// 2. Split on @ delimeter for version
	// 3. module = split[0]
	// 4. find first '/' to separate ref from file path
	// 5. version = split[1][:beforeslash]
	// 6. filepath = split[1][afterslash:]

	// Assumes no @ is anywhere else in the path
	modPathLen := len(modPath) + 1
	versionDelim := strings.Index(fullFile[modPathLen:], "@") + modPathLen
	module = fullFile[modPathLen:versionDelim]
	fileStart := strings.Index(fullFile[versionDelim:], "/") + versionDelim
	file = fullFile[fileStart+1:]
	version = fullFile[versionDelim+1 : fileStart]

	// If module is not versioned, then version tag contains a lot of extra crap on top of the commit sha
	// of the form 'v0.0.0-dateinfi-somethingelseithink-commitsha'
	// we only want the commit sha
	if strings.HasPrefix(version, "v0.0.0") {
		commitShaDelim := strings.LastIndex(version, "-")
		version = version[commitShaDelim+1:]
	}
	return
}
