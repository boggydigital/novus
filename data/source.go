package data

import (
	"errors"
	"fmt"
	"github.com/boggydigital/wits"
	"golang.org/x/net/html/atom"
	"net/url"
)

const (
	Title             = "title"
	URL               = "url"
	ContainerSelector = "container-selector"
	TextContent       = "text-content"
	ElementsAtom      = "elements-atom"
)

type Source struct {
	Id                string
	Title             string
	URL               *url.URL
	ContainerSelector string
	TextContent       string
	ElementsAtom      atom.Atom
}

func NewSource(id string, kv wits.KeyValue) (*Source, error) {
	src := &Source{
		Id: id,
	}

	for k, v := range kv {
		switch k {
		case Title:
			src.Title = v
			if src.Title == "" {
				src.Title = src.Id
			}
		case URL:
			if u, err := url.Parse(v); err != nil {
				return nil, err
			} else {
				src.URL = u
			}
		case ContainerSelector:
			src.ContainerSelector = v
		case TextContent:
			src.TextContent = v
		case ElementsAtom:
			src.ElementsAtom = atom.Lookup([]byte(v))
		default:
			return nil, errors.New("unknown source parameter " + k)
		}
	}

	return src, nil
}

func (s *Source) IsValid() error {
	if s.URL == nil {
		return fmt.Errorf("%s needs a valid %s", s.Id, URL)
	}
	if s.ContainerSelector == "" {
		return fmt.Errorf("%s requires %s", s.Id, ContainerSelector)
	}

	return nil
}
