package cli

import (
	"github.com/boggydigital/atomus"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/novus/data"
	"github.com/boggydigital/novus/rest"
	"github.com/boggydigital/pathways"
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
	NumSources int
	Added      map[string][]string
	Removed    map[string][]string
	Errors     map[string]map[string][]string
}

func PublishAtomHandler(u *url.URL) error {
	novusUrl := u.Query().Get("novus-url")
	return PublishAtom(novusUrl)
}

func PublishAtom(novusUrl string) error {

	paa := nod.NewProgress("publishing atom...")
	defer paa.End()

	properties := []string{
		data.CurrentElementsProperty,
		data.AddedElementsProperty,
		data.RemovedElementsProperty,
		data.SourceURLProperty,
	}

	// tracking errors separately from other properties to allow adding them in the future
	errorProperties := []string{
		data.GetContentErrorsProperty,
		data.MatchContentErrorsProperty,
		data.DecodeErrorsProperty,
		data.ReduceErrorsProperty,
	}

	ard, err := pathways.GetAbsDir(data.Redux)
	if err != nil {
		return paa.EndWithError(err)
	}

	properties = append(properties, errorProperties...)

	rdx, err := kvas.NewReduxReader(ard, properties...)
	if err != nil {
		return paa.EndWithError(err)
	}

	additions := make(map[string][]string)
	removals := make(map[string][]string)
	errors := make(map[string]map[string][]string)
	for _, p := range errorProperties {
		errors[p] = make(map[string][]string)
	}

	ks := keys(rdx, properties...)

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

		for _, p := range errorProperties {
			if errs, ok := rdx.GetAllValues(p, id); ok {
				for _, entry := range errs {
					errors[p][id] = append(errors[p][id], entry)
				}
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
		NumSources: len(src),
		Added:      additions,
		Removed:    removals,
		Errors:     errors,
	}
	if err := tmpl.ExecuteTemplate(sb, "atom", cl); err != nil {
		return paa.EndWithError(err)
	}
	af.SetEntry(atomEntryTitle, atomEntryAuthor, novusUrl, sb.String())

	aap, err := data.AbsAtomPath()
	if err != nil {
		return paa.EndWithError(err)
	}

	atomFile, err := os.Create(aap)
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

func keys(rdx kvas.ReadableRedux, properties ...string) map[string]interface{} {
	ks := make(map[string]interface{})
	for _, p := range properties {
		for _, id := range rdx.Keys(p) {
			ks[id] = nil
		}
	}
	return ks
}
