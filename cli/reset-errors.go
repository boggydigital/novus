package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigitl/novus/data"
	"net/url"
)

func ResetErrorsHandler(u *url.URL) error {
	return ResetErrors()
}

func ResetErrors() error {

	rca := nod.Begin("resetting errors...")
	defer rca.End()

	rdx, err := kvas.ConnectReduxAssets(
		data.AbsReduxDir(), nil,
		data.MatchContentErrorsProperty,
		data.ReduceErrorsProperty)
	if err != nil {
		return rca.EndWithError(err)
	}

	emptySet := make(map[string][]string)

	for _, id := range rdx.Keys(data.MatchContentErrorsProperty) {
		emptySet[id] = []string{}
	}

	for _, id := range rdx.Keys(data.ReduceErrorsProperty) {
		emptySet[id] = []string{}
	}

	if err := rdx.BatchReplaceValues(data.MatchContentErrorsProperty, emptySet); err != nil {
		return rca.EndWithError(err)
	}
	if err := rdx.BatchReplaceValues(data.ReduceErrorsProperty, emptySet); err != nil {
		return rca.EndWithError(err)
	}

	rca.EndWithResult("done")

	return nil
}
