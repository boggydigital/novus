package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/novus/data"
	"net/url"
)

func ResetChangesHandler(u *url.URL) error {
	return ResetChanges()
}

func ResetChanges() error {

	rca := nod.Begin("resetting changes...")
	defer rca.End()

	rdx, err := kvas.ConnectReduxAssets(
		data.AbsReduxDir(),
		data.CurrentElementsProperty,
		data.AddedElementsProperty,
		data.RemovedElementsProperty)
	if err != nil {
		return rca.EndWithError(err)
	}

	emptySet := make(map[string][]string)

	for _, id := range rdx.Keys(data.AddedElementsProperty) {
		emptySet[id] = []string{}
	}

	for _, id := range rdx.Keys(data.RemovedElementsProperty) {
		emptySet[id] = []string{}
	}

	if err := rdx.BatchReplaceValues(data.AddedElementsProperty, emptySet); err != nil {
		return rca.EndWithError(err)
	}
	if err := rdx.BatchReplaceValues(data.RemovedElementsProperty, emptySet); err != nil {
		return rca.EndWithError(err)
	}

	rca.EndWithResult("done")

	return nil
}
