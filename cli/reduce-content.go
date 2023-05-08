package cli

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/match_node"
	"github.com/boggydigital/nod"
	"github.com/boggydigitl/novus/data"
	"golang.org/x/net/html"
	"net/url"
	"strings"
)

func ReduceContentHandler(u *url.URL) error {
	return ReduceContent()
}

func ReduceContent() error {
	ra := nod.NewProgress("reducing news...")
	defer ra.End()

	rdx, err := kvas.ConnectReduxAssets(data.AbsReduxDir(), nil,
		data.CurrentElementsProperty,
		data.ReduceErrorsProperty,
		data.SourceURLProperty)
	if err != nil {
		return ra.EndWithError(err)
	}

	matchedKv, err := kvas.ConnectLocal(data.AbsMatchedContentDir(), ".html")
	if err != nil {
		return err
	}

	sources, err := loadSources()
	if err != nil {
		return err
	}

	ra.TotalInt(len(sources))

	errors := make(map[string][]string)
	urls := make(map[string][]string)

	for _, src := range sources {
		urls[src.Id] = []string{src.URL.String()}
		err = reduceSource(src, matchedKv, rdx)

		errors[src.Id] = []string{}
		if err != nil {
			errors[src.Id] = []string{err.Error()}
		}
		ra.Increment()
	}

	if err = rdx.BatchReplaceValues(data.ReduceErrorsProperty, errors); err != nil {
		return err
	}

	ra.EndWithResult("done")

	return nil
}

func reduceSource(src *data.Source, kv kvas.KeyValuesEditor, rdx kvas.ReduxAssets) error {

	if src.Query.ElementsSelector == "" {
		return nil
	}

	content, err := kv.Get(src.Id)

	if err != nil {
		return err
	}
	if content == nil {
		return fmt.Errorf("no source with id: %s", src.Id)
	}
	defer content.Close()

	doc, err := html.Parse(content)
	if err != nil {
		return err
	}

	elementsMatcher := match_node.NewSelector(src.Query.ElementsSelector)
	elements := match_node.Matches(doc, elementsMatcher, -1)

	currentElements := make([]string, 0, len(elements))

	sb := &strings.Builder{}
	for _, element := range elements {

		if src.Query.ElementReductionSelector != "" {
			reductionMatcher := match_node.NewSelector(src.Query.ElementReductionSelector)
			if reducedElement := match_node.Match(element, reductionMatcher); reducedElement != nil {
				element = reducedElement
			} else {
				continue
			}
		}

		sb.Reset()
		if err := html.Render(sb, element); err != nil {
			return err
		}
		currentElements = append(currentElements, sb.String())
	}

	if err := rdx.ReplaceValues(data.CurrentElementsProperty, src.Id, currentElements...); err != nil {
		return err
	}

	return nil
}
