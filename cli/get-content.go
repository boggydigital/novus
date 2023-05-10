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

	gna := nod.NewProgress("getting news...")
	defer gna.End()

	kv, err := kvas.ConnectLocal(data.AbsLocalContentDir(), ".html")
	if err != nil {
		return gna.EndWithError(err)
	}

	sources, err := data.LoadSources()
	if err != nil {
		return gna.EndWithError(err)
	}

	ids := make([]string, 0, len(sources))
	urls := make([]*url.URL, 0, len(sources))
	for _, src := range sources {
		ids = append(ids, src.Id)
		urls = append(urls, src.URL)
	}

	indexSetter := kvas_dolo.NewIndexSetter(kv, ids...)

	gna.TotalInt(len(sources))

	hc := http.DefaultClient
	hosts := make([]string, 0)

	if cookieFile, err := os.Open(data.AbsCookiesPath()); err == nil {
		defer cookieFile.Close()
		skv, err := wits.ReadSectionKeyValue(cookieFile)
		if err != nil {
			return gna.EndWithError(err)
		}
		for host := range skv {
			hosts = append(hosts, host)
		}
		cj, err := coost.NewJar(cookieFile, hosts...)
		if err != nil {
			return gna.EndWithError(err)
		}
		hc = cj.NewHttpClient()
	}

	dc := dolo.NewClient(hc,
		&dolo.ClientOptions{
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.4 Safari/605.1.15",
		})

	if err := dc.GetSet(urls, indexSetter, gna); err != nil {
		return gna.EndWithError(err)
	}

	gna.EndWithResult("done")

	return nil
}
