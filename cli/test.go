package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/novus/data"
	"net/url"
)

func TestSourcesHandler(u *url.URL) error {

	resetErrors := u.Query().Has("reset-errors")

	return TestSources(resetErrors)
}

func TestSources(resetErrors bool) error {

	ta := nod.Begin("testing sources...")
	defer ta.End()

	if err := GetContent(); err != nil {
		return ta.EndWithError(err)
	}

	if err := Decode(); err != nil {
		return ta.EndWithError(err)
	}

	if err := MatchContent(0); err != nil {
		return ta.EndWithError(err)
	}

	if err := Reduce(0); err != nil {
		return ta.EndWithError(err)
	}

	errorProperties := []string{
		data.GetContentErrorsProperty,
		data.DecodeErrorsProperty,
		data.MatchContentErrorsProperty,
		data.ReduceErrorsProperty,
	}

	rdx, err := kvas.ConnectReduxAssets(data.AbsReduxDir(), errorProperties...)

	if err != nil {
		return ta.EndWithError(err)
	}

	for _, p := range errorProperties {

		erra := nod.Begin("checking for %s errors...", p)

		errs := make(map[string][]string)
		for _, id := range rdx.Keys(p) {
			if errors, ok := rdx.GetAllValues(p, id); ok && len(errors) > 0 {
				errs[id] = errors
			}
		}

		if len(errs) > 0 {
			erra.EndWithSummary(p+" errors", errs)
		} else {
			erra.EndWithResult("all good")
		}
	}

	if resetErrors {
		if err := ResetErrors(); err != nil {
			return ta.EndWithError(err)
		}
	}

	ta.EndWithResult("done")

	return nil
}
