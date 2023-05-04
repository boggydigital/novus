package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/wits"
	"github.com/boggydigitl/novus/data"
	"os"
)

type sourceContentDelegate = func(source *data.Source, kv kvas.KeyValuesEditor) error

func forEachSource(scd sourceContentDelegate, tpw nod.TotalProgressWriter) error {
	kv, err := kvas.ConnectLocal(data.AbsLocalContentDir(), ".html")
	if err != nil {
		return err
	}

	sources, err := loadSources()
	if err != nil {
		return err
	}

	tpw.TotalInt(len(sources))

	for _, src := range sources {
		if err := scd(src, kv); err != nil {
			return err
		}
		tpw.Increment()
	}

	return nil
}

func loadSources() ([]*data.Source, error) {
	sf, err := os.Open(data.AbsSourcePath())
	if err != nil {
		return nil, err
	}
	defer sf.Close()

	skv, err := wits.ReadSectionKeyValue(sf)
	if err != nil {
		return nil, err
	}

	sources := make([]*data.Source, 0, len(skv))

	for id, kv := range skv {
		src, err := data.NewSource(id, kv)
		if err != nil {
			return nil, err
		}

		sources = append(sources, src)
	}

	return sources, nil
}
