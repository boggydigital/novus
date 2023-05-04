package data

import (
	"errors"
	"fmt"
	"github.com/boggydigital/wits"
)

const (
	WikipediaEnTableStudioAlbums = "wikipedia-en-table-studio-albums"
)

type QuerySelectors struct {
	ContainerSelector        string
	TextContent              string
	ElementsSelector         string
	ElementReductionSelector string
}

func NewQuerySelectors(kv wits.KeyValue) *QuerySelectors {
	qs := &QuerySelectors{}

	for k, v := range kv {
		switch k {
		case ContainerSelector:
			qs.ContainerSelector = v
		case TextContent:
			qs.TextContent = v
		case ElementsSelector:
			qs.ElementsSelector = v
		case ElementReductionSelector:
			qs.ElementReductionSelector = v
		}
	}

	return qs
}

func NewQuerySelectorsRecipe(recipe string) (*QuerySelectors, error) {
	switch recipe {
	case WikipediaEnTableStudioAlbums:
		return &QuerySelectors{
			ContainerSelector:        "table.wikitable",
			TextContent:              "List of studio albums",
			ElementsSelector:         "tr",
			ElementReductionSelector: "i",
		}, nil
	default:
		return nil, errors.New("unknown recipe " + recipe)
	}
}

func (qs *QuerySelectors) IsValid() error {
	if qs.ContainerSelector == "" {
		return fmt.Errorf("%s is required", ContainerSelector)
	}
	if qs.ElementsSelector == "" && qs.ElementReductionSelector != "" {
		return fmt.Errorf("%s is required for %s", ElementsSelector, ElementReductionSelector)
	}

	return nil
}
