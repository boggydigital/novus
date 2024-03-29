package data

import (
	"fmt"
	"github.com/boggydigital/wits"
	"net/url"
)

const (
	Id                       = "id"
	URL                      = "url"
	Recipe                   = "recipe"
	Encoding                 = "encoding"
	Title                    = "title"
	Category                 = "category"
	ContainerSelector        = "container-selector"
	TextContent              = "text-content"
	ElementsSelector         = "elements-selector"
	ElementReductionSelector = "element-reduction-selector"
	ElementAttribute         = "element-attribute"
)

type Source struct {
	Id       string
	URL      *url.URL
	Recipe   string
	Encoding string
	Title    string
	Category string
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
		case Title:
			src.Title = v
		case Category:
			src.Category = v
		case Recipe:
			src.Recipe = v
		case ContainerSelector:
			fallthrough
		case TextContent:
			fallthrough
		case ElementsSelector:
			fallthrough
		case ElementReductionSelector:
			fallthrough
		case ElementAttribute:
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
