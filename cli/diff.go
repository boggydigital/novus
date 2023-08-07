package cli

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigitl/novus/data"
	"golang.org/x/exp/slices"
	"net/url"
)

func DiffHandler(u *url.URL) error {
	return Diff()
}

func Diff() error {

	uca := nod.NewProgress("detecting changes...")
	defer uca.End()

	rdx, err := kvas.ConnectReduxAssets(
		data.AbsReduxDir(),
		data.CurrentElementsProperty,
		data.PreviousElementsProperty,
		data.AddedElementsProperty,
		data.RemovedElementsProperty)
	if err != nil {
		return uca.EndWithError(err)
	}

	currentValues := make(map[string][]string)
	addedValues := make(map[string][]string)
	removedValues := make(map[string][]string)

	for _, id := range rdx.Keys(data.CurrentElementsProperty) {

		addedValues[id] = make([]string, 0)
		removedValues[id] = make([]string, 0)

		currentValues[id], _ = rdx.GetAllValues(data.CurrentElementsProperty, id)
		previousValues, _ := rdx.GetAllValues(data.PreviousElementsProperty, id)

		for _, cv := range currentValues[id] {
			if !slices.Contains(previousValues, cv) {
				addedValues[id] = append(addedValues[id], cv)
			}
		}

		for _, pv := range previousValues {
			if !slices.Contains(currentValues[id], pv) {
				removedValues[id] = append(removedValues[id], pv)
			}
		}

	}

	if hasChanges(addedValues) || hasChanges(removedValues) {

		brv := rdx.BatchReplaceValues

		if err := brv(data.AddedElementsProperty, addedValues); err != nil {
			return uca.EndWithError(err)
		}
		if err := brv(data.RemovedElementsProperty, removedValues); err != nil {
			return uca.EndWithError(err)
		}

		// resetting previous elements - so that next diff call with the same current elements won't write anything
		if err := brv(data.PreviousElementsProperty, currentValues); err != nil {
			return uca.EndWithError(err)
		}

		uca.EndWithResult("done with changes")
	} else {
		uca.EndWithResult("no changes")
	}

	return nil
}

func hasChanges(kv map[string][]string) bool {
	for _, v := range kv {
		if len(v) > 0 {
			return true
		}
	}
	return false
}
