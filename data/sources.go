package data

import (
	"github.com/boggydigital/wits"
	"os"
)

func LoadSources() ([]*Source, error) {
	sf, err := os.Open(AbsSourcePath())
	if err != nil {
		return nil, err
	}
	defer sf.Close()

	skv, err := wits.ReadSectionKeyValue(sf)
	if err != nil {
		return nil, err
	}

	sources := make([]*Source, 0, len(skv))

	for id, kv := range skv {
		src, err := newSource(id, kv)
		if err != nil {
			return nil, err
		}

		sources = append(sources, src)
	}

	return sources, nil
}

func SourcesIds(src ...*Source) []string {
	ids := make([]string, 0, len(src))
	for _, s := range src {
		ids = append(ids, s.Id)
	}
	return ids
}

func SourceById(id string, src ...*Source) *Source {
	for _, s := range src {
		if s.Id == id {
			return s
		}
	}
	return nil
}
