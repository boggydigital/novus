package data

import "path/filepath"

const (
	localContentDirectory = "local-content"
	reduxDirectory        = "_redux"
	sourcesFilename       = "sources.txt"
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

func AbsLocalContentDir() string {
	return filepath.Join(absRootDir, localContentDirectory)
}

func AbsReduxDir() string {
	return filepath.Join(absRootDir, reduxDirectory)
}
