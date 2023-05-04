package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/match_node"
	"github.com/boggydigital/nod"
	"github.com/boggydigitl/novus/data"
	"golang.org/x/net/html"
	"net/url"
	"strings"
)

type reduxSetter struct {
	rdx kvas.ReduxValues
}

func ReduceHandler(u *url.URL) error {
	return Reduce()
}

func Reduce() error {
	ra := nod.NewProgress("reducing news...")
	defer ra.End()

	rdx, err := kvas.ConnectRedux(data.AbsReduxDir(), data.CurrentElementsProperty)
	if err != nil {
		return ra.EndWithError(err)
	}

	rs := &reduxSetter{rdx: rdx}

	if err := forEachSource(rs.reduceSource, ra); err != nil {
		return ra.EndWithError(err)
	}

	ra.EndWithResult("done")

	return nil
}

func (rs *reduxSetter) reduceSource(src *data.Source, kv kvas.KeyValuesEditor) error {

	if src.ElementsAtom == 0 {
		return nil
	}

	content, err := kv.Get(src.Id)
	if err != nil {
		return err
	}
	defer content.Close()

	doc, err := html.Parse(content)
	if err != nil {
		return err
	}

	atomMatcher := match_node.NewAtom(src.ElementsAtom)
	elements := match_node.Matches(doc, atomMatcher, -1)

	currentElements := make([]string, 0, len(elements))

	sb := &strings.Builder{}
	for _, e := range elements {
		sb.Reset()
		if err := html.Render(sb, e); err != nil {
			return err
		}
		currentElements = append(currentElements, sb.String())
	}

	if err := rs.rdx.ReplaceValues(src.Id, currentElements...); err != nil {
		return err
	}

	return nil
}
