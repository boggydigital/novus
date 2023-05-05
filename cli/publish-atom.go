package cli

import (
	"fmt"
	"github.com/boggydigital/atomus"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigitl/novus/data"
	"net/url"
	"os"
	"strings"
)

const (
	atomFeedTitle   = "novus sync updates"
	atomEntryTitle  = "Sync results "
	atomEntryAuthor = "Contemporalis"
)

func PublishAtomHandler(u *url.URL) error {
	return PublishAtom()
}

func PublishAtom() error {

	paa := nod.NewProgress("publishing atom...")
	defer paa.End()

	rdx, err := kvas.ConnectReduxAssets(
		data.AbsReduxDir(), nil,
		data.CurrentElementsProperty,
		data.AddedElementsProperty,
		data.RemovedElementsProperty,
		data.GetNewsErrorsProperty,
		data.ReduceErrorsProperty)
	if err != nil {
		return paa.EndWithError(err)
	}

	additions := make(map[string][]string)
	removals := make(map[string][]string)
	getNewsErrors := make(map[string][]string)
	reduceErrors := make(map[string][]string)

	keys := rdx.Keys(data.CurrentElementsProperty)

	paa.TotalInt(len(keys))

	for _, id := range keys {
		if added, ok := rdx.GetAllUnchangedValues(data.AddedElementsProperty, id); ok {
			for _, entry := range added {
				additions[id] = append(additions[id], entry)
			}
		}
		if removed, ok := rdx.GetAllUnchangedValues(data.RemovedElementsProperty, id); ok {
			for _, entry := range removed {
				removals[id] = append(removals[id], entry)
			}
		}
		if gnes, ok := rdx.GetAllUnchangedValues(data.GetNewsErrorsProperty, id); ok {
			for _, entry := range gnes {
				getNewsErrors[id] = append(getNewsErrors[id], entry)
			}
		}
		if res, ok := rdx.GetAllUnchangedValues(data.ReduceErrorsProperty, id); ok {
			for _, entry := range res {
				reduceErrors[id] = append(reduceErrors[id], entry)
			}
		}

		paa.Increment()
	}

	af := atomus.NewFeed(atomFeedTitle, "")
	content := changelogContent(additions, removals, getNewsErrors, reduceErrors, len(keys))
	af.SetEntry(atomEntryTitle, atomEntryAuthor, content)

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

func changelogContent(add, rem, gnes, res map[string][]string, n int) string {
	sb := &strings.Builder{}

	if len(add) > 0 {
		sb.WriteString("<h1>Added</h1>")
		for id, entries := range add {
			sb.WriteString("<h2>" + id + "</h2>")
			sb.WriteString("<ul>")
			for _, entry := range entries {
				sb.WriteString("<li>" + entry + "</li>")
			}
			sb.WriteString("</ul>")
		}
	}

	if len(rem) > 0 {
		sb.WriteString("<h1>Removed</h1>")
		for id, entries := range rem {
			sb.WriteString("<h2>" + id + "</h2>")
			sb.WriteString("<ul>")
			for _, entry := range entries {
				sb.WriteString("<li>" + entry + "</li>")
			}
			sb.WriteString("</ul>")
		}
	}

	if len(add) == 0 && len(rem) == 0 {
		sb.WriteString("<h1>No new changes since the last sync</h1>")
		sb.WriteString(fmt.Sprintf("<div>%d source(s) have been checked</div>", n))
	}

	if len(gnes) > 0 {
		sb.WriteString("<h1>get-news errors</h1>")
		for id, errors := range gnes {
			sb.WriteString("<h2>" + id + "</h2>")
			sb.WriteString("<ul>")
			for _, err := range errors {
				sb.WriteString("<li>" + err + "</li>")
			}
			sb.WriteString("</ul>")
		}
	}

	if len(res) > 0 {
		sb.WriteString("<h1>reduce errors</h1>")
		for id, errors := range res {
			sb.WriteString("<h2>" + id + "</h2>")
			sb.WriteString("<ul>")
			for _, err := range errors {
				sb.WriteString("<li>" + err + "</li>")
			}
			sb.WriteString("</ul>")
		}
	}

	return sb.String()
}
