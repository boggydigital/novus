package data

import (
	_ "embed"
	"fmt"
	"github.com/boggydigital/wits"
)

type QuerySelectors struct {
	ContainerSelector        string
	TextContent              string
	ElementsSelector         string
	ElementReductionSelector string
	ElementAttribute         string
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
		case ElementAttribute:
			qs.ElementAttribute = v
		}
	}

	return qs
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

func (qs *QuerySelectors) Override(another *QuerySelectors) {
	if another.ContainerSelector != "" {
		qs.ContainerSelector = another.ContainerSelector
	}
	if another.TextContent != "" {
		qs.TextContent = another.TextContent
	}
	if another.ElementsSelector != "" {
		qs.ElementsSelector = another.ElementsSelector
	}
	if another.ElementReductionSelector != "" {
		qs.ElementReductionSelector = another.ElementReductionSelector
	}
}
