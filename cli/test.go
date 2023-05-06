package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigitl/novus/data"
	"net/url"
)

func TestHandler(u *url.URL) error {

	resetErrors := u.Query().Has("reset-errors")

	return Test(resetErrors)
}

func Test(resetErrors bool) error {

	ta := nod.Begin("testing sources...")
	defer ta.End()

	if err := GetNews(); err != nil {
		return ta.EndWithError(err)
	}

	if err := Reduce(); err != nil {
		return ta.EndWithError(err)
	}

	rdx, err := kvas.ConnectReduxAssets(
		data.AbsReduxDir(), nil,
		data.GetNewsErrorsProperty,
		data.ReduceErrorsProperty)

	if err != nil {
		return ta.EndWithError(err)
	}

	gnea := nod.Begin("checking for get-news errors...")
	defer gnea.End()

	gnes := make(map[string][]string)
	for _, id := range rdx.Keys(data.GetNewsErrorsProperty) {
		if errors, ok := rdx.GetAllUnchangedValues(data.GetNewsErrorsProperty, id); ok && len(errors) > 0 {
			gnes[id] = errors
		}
	}

	if len(gnes) > 0 {
		gnea.EndWithSummary("get-news errors", gnes)
	} else {
		gnea.EndWithResult("all good")
	}

	rea := nod.Begin("checking for reduce errors...")
	defer rea.End()

	res := make(map[string][]string)
	for _, id := range rdx.Keys(data.ReduceErrorsProperty) {
		if errors, ok := rdx.GetAllUnchangedValues(data.ReduceErrorsProperty, id); ok && len(errors) > 0 {
			res[id] = errors
		}
	}

	if len(res) > 0 {
		rea.EndWithSummary("reduce errors", res)
	} else {
		rea.EndWithResult("all good")
	}

	if resetErrors {
		if err := ResetErrors(); err != nil {
			return ta.EndWithError(err)
		}
	}

	ta.EndWithResult("done")

	return nil
}
