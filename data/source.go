package data

import (
	"errors"
	"fmt"
	"github.com/boggydigital/wits"
	"net/url"
)

const (
	URL                      = "url"
	Recipe                   = "recipe"
	ContainerSelector        = "container-selector"
	TextContent              = "text-content"
	ElementsSelector         = "elements-selector"
	ElementReductionSelector = "element-reduction-selector"
)

type Source struct {
	Id     string
	URL    *url.URL
	Recipe string
	Query  *QuerySelectors
}

func NewSource(id string, kv wits.KeyValue) (*Source, error) {
	src := &Source{
		Id: id,
	}

	for k, v := range kv {
		switch k {
		case URL:
			if u, err := url.Parse(v); err != nil {
				return nil, err
			} else {
				src.URL = u
			}
		case Recipe:
			src.Recipe = v
		case ContainerSelector:
			fallthrough
		case TextContent:
			fallthrough
		case ElementsSelector:
			fallthrough
		case ElementReductionSelector:
			// do nothing, fill be processed later
		default:
			return nil, errors.New("unknown source parameter " + k)
		}
	}

	if src.Recipe != "" {
		var err error
		src.Query, err = NewQuerySelectorsRecipe(src.Recipe)
		if err != nil {
			return nil, err
		}
	}

	// allow overriding selectors after applying recipe
	src.Query.Override(NewQuerySelectors(kv))

	return src, nil
}

func (src *Source) IsValid() error {
	if src.URL == nil {
		return fmt.Errorf("%s needs a valid %s", src.Id, URL)
	}

	return src.Query.IsValid()
}
