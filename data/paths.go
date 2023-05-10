package data

import "path/filepath"

const (
	localContentDirectory   = "local-content"
	matchedContentDirectory = "matched-content"
	reduxDirectory          = "_redux"
	backupDirectory         = "backup"
	sourcesFilename         = "sources.txt"
	atomFilename            = "atom.xml"
	cookiesFilename         = "cookies.txt"
)

var absRootDir string

func ChRoot(d string) {
	absRootDir = d
}

func Pwd() string {
	return absRootDir
}

func AbsBackupDir() string {
	return filepath.Join(absRootDir, backupDirectory)
}

func AbsSourcePath() string {
	return filepath.Join(absRootDir, sourcesFilename)
}

func AbsAtomPath() string {
	return filepath.Join(absRootDir, atomFilename)
}

func AbsCookiesPath() string {
	return filepath.Join(absRootDir, cookiesFilename)
}

func AbsLocalContentDir() string {
	return filepath.Join(absRootDir, localContentDirectory)
}

func AbsMatchedContentDir() string {
	return filepath.Join(absRootDir, matchedContentDirectory)
}

func AbsReduxDir() string {
	return filepath.Join(absRootDir, reduxDirectory)
}
