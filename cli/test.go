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

	if err := GetContent(); err != nil {
		return ta.EndWithError(err)
	}

	if err := MatchContent(0); err != nil {
		return ta.EndWithError(err)
	}

	if err := ReduceContent(0); err != nil {
		return ta.EndWithError(err)
	}

	rdx, err := kvas.ConnectReduxAssets(
		data.AbsReduxDir(), nil,
		data.MatchContentErrorsProperty,
		data.ReduceErrorsProperty)

	if err != nil {
		return ta.EndWithError(err)
	}

	gnea := nod.Begin("checking for match-content errors...")
	defer gnea.End()

	mces := make(map[string][]string)
	for _, id := range rdx.Keys(data.MatchContentErrorsProperty) {
		if errors, ok := rdx.GetAllUnchangedValues(data.MatchContentErrorsProperty, id); ok && len(errors) > 0 {
			mces[id] = errors
		}
	}

	if len(mces) > 0 {
		gnea.EndWithSummary("match-content errors", mces)
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
