package cli

import (
	"errors"
	"fmt"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/match_node"
	"github.com/boggydigital/nod"
	"github.com/boggydigitl/novus/data"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"strings"
)

func GetNewsHandler(u *url.URL) error {
	return GetNews()
}

func GetNews() error {

	gna := nod.NewProgress("getting news...")
	defer gna.End()

	if err := forEachSource(getSource, data.GetNewsErrorsProperty, gna); err != nil {
		return gna.EndWithError(err)
	}

	gna.EndWithResult("done")

	return nil
}

func getSource(src *data.Source, kv kvas.KeyValuesEditor) error {
	gsa := nod.Begin(" getting %s...", src.Id)
	defer gsa.End()

	if err := src.IsValid(); err != nil {
		return gsa.EndWithError(err)
	}

	resp, err := http.Get(src.URL.String())
	if err != nil {
		return gsa.EndWithError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return gsa.EndWithError(errors.New(resp.Status))
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return gsa.EndWithError(err)
	}

	matches := match_node.Matches(doc, match_node.NewSelector(src.Query.ContainerSelector), -1)

	sb := &strings.Builder{}
	textContent := ""
	contentSelected := false

	for _, m := range matches {
		sb.Reset()
		if err := html.Render(sb, m); err != nil {
			return gsa.EndWithError(err)
		}

		textContent = sb.String()
		if strings.Contains(textContent, src.Query.TextContent) {
			contentSelected = true

			if err := kv.Set(src.Id, strings.NewReader(textContent)); err != nil {
				return gsa.EndWithError(err)
			}

			break
		}
	}

	if !contentSelected {
		err := fmt.Errorf("%s doesn't have elements with: '%s'", src.Id, data.TextContent)
		return gsa.EndWithError(err)
	}

	gsa.EndWithResult("done")

	return nil
}
