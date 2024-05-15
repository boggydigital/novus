package data

import (
	"github.com/boggydigital/pathways"
	"path/filepath"
)

const (
	sourcesFilename = "sources.txt"
	atomFilename    = "atom.xml"
	cookiesFilename = "cookies.txt"
)

func AbsSourcesPath() (string, error) {
	aid, err := pathways.GetAbsDir(Input)
	if err != nil {
		return "", err
	}

	return filepath.Join(aid, sourcesFilename), nil
}

func AbsCookiesPath() (string, error) {
	aid, err := pathways.GetAbsDir(Input)
	if err != nil {
		return "", err
	}

	return filepath.Join(aid, cookiesFilename), nil
}

func AbsAtomPath() (string, error) {
	aod, err := pathways.GetAbsDir(Output)
	if err != nil {
		return "", err
	}

	return filepath.Join(aod, atomFilename), nil
}
