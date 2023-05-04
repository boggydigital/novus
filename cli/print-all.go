package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigitl/novus/data"
	"net/url"
)

type reduxGetter struct {
	rdx kvas.ReduxValues
}

func PrintAllHandler(u *url.URL) error {
	return PrintAll()
}

func PrintAll() error {
	paa := nod.NewProgress("printing all...")
	defer paa.End()

	rdx, err := kvas.ConnectRedux(data.AbsReduxDir(), data.CurrentElementsProperty)
	if err != nil {
		return paa.EndWithError(err)
	}

	rs := &reduxGetter{rdx: rdx}

	if err := forEachSource(rs.printAll, paa); err != nil {
		return paa.EndWithError(err)
	}

	paa.EndWithResult("done")

	return nil
}

func (rg *reduxGetter) printAll(src *data.Source, _ kvas.KeyValuesEditor) error {

	ia := nod.Begin(src.Id)
	defer ia.End()

	values, ok := rg.rdx.GetAllValues(src.Id)
	if !ok {
		return nil
	}
	for _, val := range values {
		va := nod.Begin(" %s", val)
		va.End()
	}

	return nil
}
