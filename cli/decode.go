package cli

import (
	"errors"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/novus/data"
	"github.com/boggydigital/pathways"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"io"
	"net/url"
)

func DecodeHandler(u *url.URL) error {
	return Decode()
}

func Decode() error {

	dca := nod.NewProgress("decoding content...")
	defer dca.End()

	alcd, err := pathways.GetAbsDir(data.LocalContent)
	if err != nil {
		return dca.EndWithError(err)
	}

	ard, err := pathways.GetAbsDir(data.Redux)
	if err != nil {
		return dca.EndWithError(err)
	}

	kv, err := kvas.ConnectLocal(alcd, ".html")
	if err != nil {
		return dca.EndWithError(err)
	}

	rdx, err := kvas.NewReduxWriter(ard, data.DecodeErrorsProperty)
	if err != nil {
		return dca.EndWithError(err)
	}

	sources, err := data.LoadSources()
	if err != nil {
		return dca.EndWithError(err)
	}

	dca.TotalInt(len(sources))

	errs := make(map[string][]string)

	for _, src := range sources {

		if src.Encoding == "" {
			continue
		}

		if err := decodeContent(kv, src.Id, src.Encoding); err != nil {
			errs[src.Id] = []string{err.Error()}
		}

		dca.Increment()
	}

	if err := rdx.BatchReplaceValues(data.DecodeErrorsProperty, errs); err != nil {
		return dca.EndWithError(err)
	}

	dca.EndWithResult("done")

	return nil
}

func decodeContent(kv kvas.KeyValues, id, enc string) error {

	rc, err := kv.Get(id)
	if err != nil {
		return err
	}
	defer rc.Close()

	var dc *encoding.Decoder
	var dr io.Reader

	switch enc {
	case "win1251":
		dc = charmap.Windows1251.NewDecoder()
	default:
		return errors.New("unknown encoding " + enc)
	}

	if dc != nil {
		dr = dc.Reader(rc)
	}

	if dr != nil {
		return kv.Set(id, dr)
	}

	return errors.New("content not decoded")
}
