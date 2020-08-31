package mod

import "strings"

func ExtractVersion(path string) (newpath, ver string) {
	newpath = path
	s := strings.Split(path, "@")
	if len(s) > 1 {
		ver = s[len(s)-1]
		newpath = path[:len(path)-len(ver)-1]
	}
	return
}

func AppendVersion(path, ver string) string {
	if ver == "" {
		return path
	}
	return path + "@" + ver
}
