package cli

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/novus/data"
	"github.com/boggydigital/pathways"
	"net/url"
)

func ResetErrorsHandler(u *url.URL) error {
	return ResetErrors()
}

func ResetErrors() error {

	rca := nod.Begin("resetting errors...")
	defer rca.End()

	ard, err := pathways.GetAbsDir(data.Redux)
	if err != nil {
		return rca.EndWithError(err)
	}

	errorProperties := []string{
		data.GetContentErrorsProperty,
		data.MatchContentErrorsProperty,
		data.DecodeErrorsProperty,
		data.ReduceErrorsProperty,
	}

	rdx, err := kevlar.NewReduxWriter(ard, errorProperties...)
	if err != nil {
		return rca.EndWithError(err)
	}

	emptySet := make(map[string]map[string][]string)

	for _, p := range errorProperties {
		if emptySet[p] == nil {
			emptySet[p] = make(map[string][]string)
		}
		for _, id := range rdx.Keys(p) {
			emptySet[p][id] = []string{}
		}

		if err := rdx.BatchReplaceValues(p, emptySet[p]); err != nil {
			return rca.EndWithError(err)
		}
	}

	rca.EndWithResult("done")

	return nil
}
