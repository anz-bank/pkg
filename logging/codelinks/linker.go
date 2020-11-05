package codelinks

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

var (
	// build constants to control the behaviour of GetCodeLinker
	// to use remote links in build images, set the following constants in ldflags at build time
	//
	// github.com/anz-bank/pkg/logging/codelinks.localLinks=false
	// github.com/anz-bank/pkg/logging/codelinks.projectRoot=$(git rev-parse --show-toplevel)
	// github.com/anz-bank/pkg/logging/codelinks.repoURL=<hostURL>
	// github.com/anz-bank/pkg/logging/codelinks.version=<version>
	//
	// WARNING: GOPATH can be multivalued, so the above might not work. path should point to where
	// modules are stored
	localLinks  = "true"
	projectRoot string
	repoURL     string
	version     string

	// modpath is a build constant that must be set when building without vendor
	//
	// github.com/anz-bank/pkg/logging/codelinks.modpath=$GOPATH/pkg/mod
	modpath string
)

var (
	// local link producer needs to switch slashes if running on windows.
	// For some reason runtime.Caller returns forward slashes in file paths
	localFileFunc = func() func(string) string {
		if runtime.GOOS != "windows" {
			return func(s string) string { return s }
		}
		return func(s string) string {
			return strings.ReplaceAll(s, `/`, `\`)
		}
	}()
)

// GetCodeLinker gets the application's code linker
//
// Return value depends on the build constants. See project documentation for more info.
// Defaults use local links
//
// Errors returned here are not serious, and loggers will work without a linker, but it
// is recommended you at least log the error
func GetCodeLinker() (CodeLinker, error) {
	_, file, _, _ := runtime.Caller(0)
	return getCodeLinker(codeLinkerConfig{
		localLinks:  localLinks == "true",
		fromVendor:  strings.HasPrefix(file, projectRoot+"/vendor"),
		projectRoot: projectRoot,
		modpath:     modpath,
		repoURL:     repoURL,
		version:     version,
	})
}

type codeLinkerConfig struct {
	localLinks  bool
	fromVendor  bool
	projectRoot string
	modpath     string
	repoURL     string
	version     string
}

func getCodeLinker(cfg codeLinkerConfig) (CodeLinker, error) {
	if cfg.localLinks {
		return LocalLinker{}, nil
	}
	if cfg.fromVendor {
		return newVendorLinker(cfg.projectRoot, cfg.repoURL, cfg.version)
	}
	if cfg.modpath != "" {
		return newModLinker(cfg.projectRoot, cfg.modpath, cfg.repoURL, cfg.version)
	}
	return nil, errors.New("logging.CodeLinker: Codelink configuration insufficient to support code links")
}

// CodeLinker defines objects that can produce source code links from file and line number
type CodeLinker interface {
	// Link produces a source code link from source file and line number
	Link(file string, line int) string
}

// LocalLinker produces links to the local file system
//
// The local file system is actually the file system at build time, so is only
// useful in development where build and run occur on the same machine
type LocalLinker struct{}

// Link produces a source code link to the local file system
func (l LocalLinker) Link(file string, line int) string {
	return localFileFunc(file) + ":" + strconv.Itoa(line)
}

func createLink(module, file, version string, line int) string {
	if strings.HasPrefix(module, "github") {
		return GithubLink(module, file, version, line)
	}
	return notFound
}

// GithubLink creates a link to a github site
//
// NOTE: github.com is NOT hard coded as ghe instances are not hosted on github.com, yet this code works just fine
func GithubLink(module string, file string, ref string, line int) string {
	// A module may extend into a github repo, rather than just point directly at it
	// In this case, we need to split the module after the third path component and
	// attach the remainder to the file to construct the link
	const githubRepoPathLength = 3
	split := strings.Split(module, "/")
	if len(split) > githubRepoPathLength {
		file = strings.Join(append(split[githubRepoPathLength:], file), "/")
		module = strings.Join(split[:githubRepoPathLength], "/")
	}
	return fmt.Sprintf("https://%s/tree/%s/%s#L%d", module, ref, file, line)
}
