package cli

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/match_node"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/novus/data"
	"golang.org/x/net/html"
	"net/url"
	"strings"
)

func MatchContentHandler(u *url.URL) error {
	return MatchContent()
}

func MatchContent() error {

	mca := nod.NewProgress("matching content...")
	defer mca.End()

	localKv, err := kvas.ConnectLocal(data.AbsLocalContentDir(), ".html")
	if err != nil {
		return mca.EndWithError(err)
	}

	matchedKv, err := kvas.ConnectLocal(data.AbsMatchedContentDir(), ".html")
	if err != nil {
		return mca.EndWithError(err)
	}

	sources, err := data.LoadSources()
	if err != nil {
		return mca.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(data.AbsReduxDir(), data.MatchContentErrorsProperty)
	if err != nil {
		return mca.EndWithError(err)
	}

	errors := make(map[string][]string)

	ids := data.SourcesIds(sources...)

	mca.TotalInt(len(ids))

	for _, id := range ids {

		src := data.SourceById(id, sources...)

		errors[id] = make([]string, 0)
		if err := matchSource(src, localKv, matchedKv); err != nil {
			errors[id] = []string{err.Error()}
			mca.Increment()
			continue
		}

		mca.Increment()
	}

	if len(errors) > 0 {
		if err := rdx.BatchReplaceValues(data.MatchContentErrorsProperty, errors); err != nil {
			return mca.EndWithError(err)
		}
	}

	mca.EndWithResult("done")

	return nil
}

func matchSource(src *data.Source, localKv, matchedKv kvas.KeyValues) error {
	localContent, err := localKv.Get(src.Id)
	if err != nil {
		return err
	}
	defer localContent.Close()

	doc, err := html.Parse(localContent)
	if err != nil {
		return err
	}

	matches := match_node.Matches(doc, match_node.NewSelector(src.Query.ContainerSelector), -1)
	if len(matches) == 0 {
		err := fmt.Errorf("%s %s selected nothing", src.Id, data.ContainerSelector)
		return err
	}

	sb := &strings.Builder{}
	textContent := ""
	contentSelected := false

	for _, m := range matches {
		sb.Reset()
		if err := html.Render(sb, m); err != nil {
			return err
		}

		textContent = sb.String()
		if strings.Contains(textContent, src.Query.TextContent) {
			contentSelected = true

			if err := matchedKv.Set(src.Id, strings.NewReader(textContent)); err != nil {
				return err
			}

			break
		}
	}

	if !contentSelected {
		err := fmt.Errorf("%s %s selected nothing", src.Id, data.TextContent)
		return err
	}

	return nil
}
