package mod

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/afero"
)

type Module struct {
	// The name of module joined by forward slash(/). e.g. "github.com/anz-bank/foo"
	Name string
	// The absolute path to the module.
	// e.g. "/Users/username/go/pkg/mod/github.com/anz-bank/foo@v1.1.0" on Linux and macOS
	//      "C:\Users\username\go\pkg\mod\github.com\anz-bank\foo@v1.1.0" on Windows
	Dir string
	// The version of the module. e.g. "v1.1.0"
	Version string
}

type Modules []*Module

var modules Modules
var manager DependencyManager = &goModules{}
var mode ModeType = GoModulesMode

type ModeType string

const (
	GitHubMode    ModeType = "github"
	GoModulesMode ModeType = "go modules"
)
const MasterBranch = "master"

type DependencyManager interface {
	// Download external dependency to local directory
	Get(filename, ver string, m *Modules) (*Module, error)
	// Find dependency in m *Modules
	Find(filename, ver string, m *Modules) *Module
	// Load local cache into m *Modules
	Load(m *Modules) error
}

func (m *Modules) Add(v *Module) {
	*m = append(*m, v)
}

func (m *Modules) Len() int {
	return len(*m)
}

var configMutex = &sync.RWMutex{}

func Config(m ModeType, goModopt GoModulesOptions, githubOpt GitHubOptions) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	mode = m
	switch mode {
	case GitHubMode:
		gh := &githubMgr{}
		if err := gh.Init(githubOpt); err != nil {
			return err
		}
		manager = gh
	case GoModulesMode:
		gm := &goModules{}
		if err := gm.Init(goModopt); err != nil {
			return err
		}
		manager = gm
	default:
		return fmt.Errorf("unknown mode type %s", mode)
	}
	return nil
}

func Retrieve(name string, ver string) (*Module, error) {
	if modules.Len() == 0 {
		if err := manager.Load(&modules); err != nil {
			return nil, fmt.Errorf("error loading modules: %s", err.Error())
		}
	}

	if ver != MasterBranch || (mode == GitHubMode && ver != "") {
		mod := manager.Find(name, ver, &modules)
		if mod != nil {
			return mod, nil
		}
	}

	return manager.Get(name, ver, &modules)
}

func hasPathPrefix(prefix, s string) bool {
	prefix = filepath.Clean(prefix)
	s = filepath.Clean(s)

	if len(s) > len(prefix) {
		return s[len(prefix)] == filepath.Separator && s[:len(prefix)] == prefix
	}

	return s == prefix
}

func FileExists(fs afero.Fs, filename string, isDir bool) bool {
	info, err := fs.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if isDir {
		return info.IsDir()
	}
	return !info.IsDir()
}
