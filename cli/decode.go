package cli

import (
	"errors"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigitl/novus/data"
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

	kv, err := kvas.ConnectLocal(data.AbsLocalContentDir(), ".html")
	if err != nil {
		return dca.EndWithError(err)
	}

	sources, err := data.LoadSources()
	if err != nil {
		return dca.EndWithError(err)
	}

	dca.TotalInt(len(sources))

	for _, src := range sources {

		if src.Query.Encoding == "" {
			continue
		}

		if err := decodeContent(kv, src.Id, src.Query.Encoding); err != nil {
			return dca.EndWithError(err)
		}

		dca.Increment()
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
