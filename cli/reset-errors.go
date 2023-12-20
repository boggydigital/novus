package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/novus/data"
	"net/url"
)

func ResetErrorsHandler(u *url.URL) error {
	return ResetErrors()
}

func ResetErrors() error {

	rca := nod.Begin("resetting errors...")
	defer rca.End()

	errorProperties := []string{
		data.GetContentErrorsProperty,
		data.MatchContentErrorsProperty,
		data.DecodeErrorsProperty,
		data.ReduceErrorsProperty,
	}

	rdx, err := kvas.ReduxWriter(data.AbsReduxDir(), errorProperties...)
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
