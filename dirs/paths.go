package dirs

import "path/filepath"

const (
	sourcesFilename = "sources.txt"
)

var absRootDir string

func ChRoot(d string) {
	absRootDir = d
}

func Pwd() string {
	return absRootDir
}

func AbsSourcePath() string {
	return filepath.Join(absRootDir, sourcesFilename)
}
