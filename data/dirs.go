package data

import "github.com/boggydigital/pathways"

const DefaultNovusRootDir = "/usr/share/novus"

const (
	Backups        pathways.AbsDir = "backups"
	Input          pathways.AbsDir = "input"
	LocalContent   pathways.AbsDir = "local-content"
	MatchedContent pathways.AbsDir = "matched-content"
	Output         pathways.AbsDir = "output"
	Redux          pathways.AbsDir = "redux"
)

func AllAbsDirs() []pathways.AbsDir {
	return []pathways.AbsDir{
		Backups,
		Input,
		LocalContent,
		MatchedContent,
		Output,
		Redux,
	}
}
