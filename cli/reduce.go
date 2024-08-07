package cli

import (
	"fmt"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/match_node"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/novus/data"
	"github.com/boggydigital/pathways"
	"golang.org/x/exp/slices"
	"golang.org/x/net/html"
	"net/url"
	"strings"
)

func ReduceHandler(u *url.URL) error {
	return Reduce(0)
}

func Reduce(since int64) error {
	rca := nod.NewProgress("reducing content...")
	defer rca.End()

	amcd, err := pathways.GetAbsDir(data.MatchedContent)
	if err != nil {
		return rca.EndWithError(err)
	}

	ard, err := pathways.GetAbsDir(data.Redux)
	if err != nil {
		return rca.EndWithError(err)
	}

	rdx, err := kevlar.NewReduxWriter(ard,
		data.CurrentElementsProperty,
		data.ReduceErrorsProperty,
		data.SourceURLProperty)
	if err != nil {
		return rca.EndWithError(err)
	}

	matchedKv, err := kevlar.NewKeyValues(amcd, kevlar.HtmlExt)
	if err != nil {
		return rca.EndWithError(err)
	}

	sources, err := data.LoadSources()
	if err != nil {
		return rca.EndWithError(err)
	}

	rca.TotalInt(len(sources))

	errors := make(map[string][]string)
	urls := make(map[string][]string)

	var ids []string
	if since > 0 {
		ids, err = matchedKv.CreatedOrUpdatedAfter(since)
		if err != nil {
			return rca.EndWithError(err)
		}
	} else {
		ids = data.SourcesIds(sources...)
	}

	for _, id := range ids {
		src := data.SourceById(id, sources...)

		urls[id] = []string{src.URL.String()}
		err = reduceSource(src, matchedKv, rdx)

		errors[id] = []string{}
		if err != nil {
			errors[id] = []string{err.Error()}
		}
		rca.Increment()
	}

	if err = rdx.BatchReplaceValues(data.SourceURLProperty, urls); err != nil {
		return rca.EndWithError(err)
	}
	if err = rdx.BatchReplaceValues(data.ReduceErrorsProperty, errors); err != nil {
		return rca.EndWithError(err)
	}

	rca.EndWithResult("done")

	return nil
}

func reduceSource(src *data.Source, kv kevlar.KeyValues, rdx kevlar.WriteableRedux) error {

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

	for _, element := range elements {

		if src.Query.ElementReductionSelector != "" {
			reductionMatcher := match_node.NewSelector(src.Query.ElementReductionSelector)
			if reducedElement := match_node.Match(element, reductionMatcher); reducedElement != nil {
				element = reducedElement
			} else {
				continue
			}
		}

		tc := ""
		if src.Query.ElementAttribute == "" {
			tc = match_node.TextContent(element)
			// clean up extra white-space characters
			tc = strings.Replace(tc, "\n", "", -1)
			tc = strings.Replace(tc, "\t", "", -1)
		} else {
			tc = match_node.AttrVal(element, src.Query.ElementAttribute)
			if tc == "" {
				return fmt.Errorf("%s[%s] attribute has no value", src.Id, src.Query.ElementAttribute)
			}
		}

		if tc != "" && !slices.Contains(currentElements, tc) {
			currentElements = append(currentElements, tc)
		}
	}

	if err := rdx.ReplaceValues(data.CurrentElementsProperty, src.Id, currentElements...); err != nil {
		return err
	}

	return nil
}
