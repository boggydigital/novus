package cli

import (
	"errors"
	"github.com/boggydigital/coost"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/novus/data"
	"github.com/boggydigital/pathways"
	"net/http"
	"net/url"
)

func GetContentHandler(u *url.URL) error {
	return GetContent()
}

func GetContent() error {

	gca := nod.NewProgress("getting content...")
	defer gca.End()

	alcd, err := pathways.GetAbsDir(data.LocalContent)
	if err != nil {
		return gca.EndWithError(err)
	}

	ard, err := pathways.GetAbsDir(data.Redux)
	if err != nil {
		return gca.EndWithError(err)
	}

	kv, err := kevlar.NewKeyValues(alcd, kevlar.HtmlExt)
	if err != nil {
		return gca.EndWithError(err)
	}

	rdx, err := kevlar.NewReduxWriter(ard, data.GetContentErrorsProperty)
	if err != nil {
		return gca.EndWithError(err)
	}

	sources, err := data.LoadSources()
	if err != nil {
		return gca.EndWithError(err)
	}

	gca.TotalInt(len(sources))

	acp, err := data.AbsCookiesPath()
	if err != nil {
		return gca.EndWithError(err)
	}

	hc, err := coost.NewHttpClientFromFile(acp)
	if err != nil {
		return gca.EndWithError(err)
	}

	getContentErrors := make(map[string][]string)

	for _, src := range sources {
		if err := getSource(src, hc, kv); err != nil {
			getContentErrors[src.Id] = []string{err.Error()}
			continue
		}
		gca.Increment()
	}

	if err := rdx.BatchReplaceValues(data.GetContentErrorsProperty, getContentErrors); err != nil {
		return gca.EndWithError(err)
	}

	gca.EndWithResult("done")

	return nil
}

func getSource(src *data.Source, hc *http.Client, kv kevlar.KeyValues) error {

	resp, err := hc.Get(src.URL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New(resp.Status)
	}

	if err := kv.Set(src.Id, resp.Body); err != nil {
		return err
	}

	return nil
}
