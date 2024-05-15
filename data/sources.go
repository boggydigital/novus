package data

import (
	"github.com/boggydigital/wits"
	"os"
)

func LoadSources() ([]*Source, error) {

	asp, err := AbsSourcesPath()
	if err != nil {
		return nil, err
	}

	sf, err := os.Open(asp)
	if err != nil {
		return nil, err
	}
	defer sf.Close()

	skv, err := wits.ReadSectionKeyValue(sf)
	if err != nil {
		return nil, err
	}

	recipes, err := LoadRecipes()
	if err != nil {
		return nil, err
	}

	sources := make([]*Source, 0, len(skv))

	for id, kv := range skv {
		src, err := newSource(id, kv, recipes)
		if err != nil {
			return nil, err
		}

		sources = append(sources, src)
	}

	return sources, validSources(sources...)
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

func validSources(src ...*Source) error {
	for _, s := range src {
		if err := s.IsValid(); err != nil {
			return err
		}
	}
	return nil
}
