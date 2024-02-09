package cli

import (
	"errors"
	"github.com/boggydigital/coost"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/novus/data"
	"net/http"
	"net/url"
)

func GetContentHandler(u *url.URL) error {
	return GetContent()
}

func GetContent() error {

	gca := nod.NewProgress("getting content...")
	defer gca.End()

	kv, err := kvas.ConnectLocal(data.AbsLocalContentDir(), ".html")
	if err != nil {
		return gca.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(data.AbsReduxDir(), data.GetContentErrorsProperty)
	if err != nil {
		return gca.EndWithError(err)
	}

	sources, err := data.LoadSources()
	if err != nil {
		return gca.EndWithError(err)
	}

	gca.TotalInt(len(sources))

	cj, err := coost.NewJar(data.AbsCookiesPath())
	if err != nil {
		return gca.EndWithError(err)
	}

	hc := cj.NewHttpClient()
	getContentErrors := make(map[string][]string)

	for _, src := range sources {
		if err := getSource(src, hc, kv); err != nil {
			getContentErrors[src.Id] = []string{err.Error()}
			continue
		}
		if err := cj.Store(data.AbsCookiesPath()); err != nil {
			return gca.EndWithError(err)
		}
		gca.Increment()
	}

	if err := rdx.BatchReplaceValues(data.GetContentErrorsProperty, getContentErrors); err != nil {
		return gca.EndWithError(err)
	}

	gca.EndWithResult("done")

	return nil
}

func getSource(src *data.Source, hc *http.Client, kv kvas.KeyValues) error {

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
