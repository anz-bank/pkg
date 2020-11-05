package codelinks

import (
	"strings"
)

// VendorLinker implements CodeLinker for projects that build from vendor
type vendorLinker struct {
	ProjectRoot string
	projectLink func(string, int) string
}

func newVendorLinker(projectRoot, repoURL, projectVersion string) (*vendorLinker, error) {
	var projectLink func(string, int) string
	if strings.Contains(repoURL, "github") {
		repo := githubRepoFromURL(repoURL)
		projectLink = func(file string, line int) string {
			return GithubLink(repo, file, projectVersion, line)
		}
	} else {
		return nil, ErrUnsupportedRemoteRepository
	}
	return &vendorLinker{
		ProjectRoot: projectRoot,
		projectLink: projectLink,
	}, nil
}

// Link implements CodeLinker.Link
func (v *vendorLinker) Link(file string, line int) string {
	file = file[len(v.ProjectRoot)+1:]
	return v.projectLink(file, line)
}
