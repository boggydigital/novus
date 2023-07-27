package data

import (
	"fmt"
	"github.com/boggydigital/wits"
	"net/url"
)

const (
	URL                      = "url"
	Recipe                   = "recipe"
	Encoding                 = "encoding"
	ContainerSelector        = "container-selector"
	TextContent              = "text-content"
	ElementsSelector         = "elements-selector"
	ElementReductionSelector = "element-reduction-selector"
)

type Source struct {
	Id       string
	URL      *url.URL
	Recipe   string
	Encoding string
	Query    *QuerySelectors
}

func newSource(id string, kv wits.KeyValue, recipes wits.SectionKeyValue) (*Source, error) {
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
		case Encoding:
			src.Encoding = v
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
			return nil, fmt.Errorf("%s has an unknown source parameter %s", id, k)
		}
	}

	if src.Recipe != "" {
		if rcp, ok := recipes[src.Recipe]; ok {
			src.Query = NewQuerySelectors(rcp)
		}
	}

	if src.Query == nil {
		src.Query = &QuerySelectors{}
	}

	// allow overriding selectors after applying recipe
	if qs := NewQuerySelectors(kv); qs != nil {
		src.Query.Override(qs)
	}

	return src, nil
}

func (src *Source) IsValid() error {
	if src.URL == nil {
		return fmt.Errorf("%s needs a valid %s", src.Id, URL)
	}

	return src.Query.IsValid()
}
