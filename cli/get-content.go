package cli

import (
	"github.com/boggydigital/coost"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/kvas_dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/wits"
	"github.com/boggydigitl/novus/data"
	"net/http"
	"net/url"
	"os"
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

	rdx, err := kvas.ConnectRedux(data.AbsReduxDir(), data.GetContentErrorsProperty)
	if err != nil {
		return gca.EndWithError(err)
	}

	sources, err := data.LoadSources()
	if err != nil {
		return gca.EndWithError(err)
	}

	ids := make([]string, 0, len(sources))
	urls := make([]*url.URL, 0, len(sources))
	for _, src := range sources {
		ids = append(ids, src.Id)
		urls = append(urls, src.URL)
	}

	indexSetter := kvas_dolo.NewIndexSetter(kv, ids...)

	gca.TotalInt(len(sources))

	hc := http.DefaultClient
	hosts := make([]string, 0)

	if cookieFile, err := os.Open(data.AbsCookiesPath()); err == nil {
		defer cookieFile.Close()
		skv, err := wits.ReadSectionKeyValue(cookieFile)
		if err != nil {
			return gca.EndWithError(err)
		}
		for host := range skv {
			hosts = append(hosts, host)
		}
		cj, err := coost.NewJar(cookieFile, hosts...)
		if err != nil {
			return gca.EndWithError(err)
		}
		hc = cj.NewHttpClient()
	}

	dc := dolo.NewClient(hc, dolo.Defaults())

	if errs := dc.GetSet(urls, indexSetter, gca); len(errs) > 0 {
		sourceErrors := make(map[string][]string, len(errs))
		for i, src := range sources {
			if err, ok := errs[i]; ok && err != nil {
				sourceErrors[src.Id] = []string{err.Error()}
			}
		}
		if err := rdx.BatchReplaceValues(sourceErrors); err != nil {
			return gca.EndWithError(err)
		}
	}

	gca.EndWithResult("done")

	return nil
}
