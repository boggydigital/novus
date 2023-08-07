package cli

import (
	"github.com/boggydigital/atomus"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigitl/novus/data"
	"github.com/boggydigitl/novus/rest"
	"net/url"
	"os"
	"strings"
)

const (
	atomFeedTitle   = "novus sync updates"
	atomEntryTitle  = "Sync results"
	atomEntryAuthor = "Contemporalis"
)

type Changelog struct {
	NumSources          int
	Added               map[string][]string
	Removed             map[string][]string
	MatchContentErrors  map[string][]string
	ReduceContentErrors map[string][]string
}

func PublishAtomHandler(u *url.URL) error {
	novusUrl := u.Query().Get("novus-url")
	return PublishAtom(novusUrl)
}

func PublishAtom(novusUrl string) error {

	paa := nod.NewProgress("publishing atom...")
	defer paa.End()

	rdx, err := kvas.ConnectReduxAssets(
		data.AbsReduxDir(),
		data.CurrentElementsProperty,
		data.AddedElementsProperty,
		data.RemovedElementsProperty,
		data.MatchContentErrorsProperty,
		data.ReduceErrorsProperty,
		data.SourceURLProperty)
	if err != nil {
		return paa.EndWithError(err)
	}

	additions := make(map[string][]string)
	removals := make(map[string][]string)
	matchContentErrors := make(map[string][]string)
	reduceErrors := make(map[string][]string)

	ks := keys(rdx,
		data.CurrentElementsProperty,
		data.AddedElementsProperty,
		data.RemovedElementsProperty,
		data.MatchContentErrorsProperty,
		data.ReduceErrorsProperty)

	paa.TotalInt(len(ks))

	for id := range ks {
		host := ""
		if su, ok := rdx.GetFirstVal(data.SourceURLProperty, id); ok {
			host = rest.Host(su)
		}

		if added, ok := rdx.GetAllValues(data.AddedElementsProperty, id); ok {
			for _, entry := range added {
				additions[id] = append(additions[id], rest.AbsHref(entry, host))
			}
		}
		if removed, ok := rdx.GetAllValues(data.RemovedElementsProperty, id); ok {
			for _, entry := range removed {
				removals[id] = append(removals[id], rest.AbsHref(entry, host))
			}
		}
		if mces, ok := rdx.GetAllValues(data.MatchContentErrorsProperty, id); ok {
			for _, entry := range mces {
				matchContentErrors[id] = append(matchContentErrors[id], entry)
			}
		}
		if res, ok := rdx.GetAllValues(data.ReduceErrorsProperty, id); ok {
			for _, entry := range res {
				reduceErrors[id] = append(reduceErrors[id], entry)
			}
		}

		paa.Increment()
	}

	src, err := data.LoadSources()
	if err != nil {
		return paa.EndWithError(err)
	}

	af := atomus.NewFeed(atomFeedTitle, novusUrl)
	sb := &strings.Builder{}
	cl := &Changelog{
		NumSources:          len(src),
		Added:               additions,
		Removed:             removals,
		MatchContentErrors:  matchContentErrors,
		ReduceContentErrors: reduceErrors,
	}
	if err := tmpl.ExecuteTemplate(sb, "atom", cl); err != nil {
		return paa.EndWithError(err)
	}
	af.SetEntry(atomEntryTitle, atomEntryAuthor, novusUrl, sb.String())

	atomFile, err := os.Create(data.AbsAtomPath())
	if err != nil {
		return paa.EndWithError(err)
	}
	defer atomFile.Close()

	if err := af.Encode(atomFile); err != nil {
		return paa.EndWithError(err)
	}

	paa.EndWithResult("done")

	return nil
}

func keys(rdx kvas.ReduxAssets, properties ...string) map[string]interface{} {
	ks := make(map[string]interface{})
	for _, p := range properties {
		for _, id := range rdx.Keys(p) {
			ks[id] = nil
		}
	}
	return ks
}
