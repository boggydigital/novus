package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/novus/data"
	"net/http"
)

func GetSource(w http.ResponseWriter, r *http.Request) {

	// GET /source?id

	id := r.URL.Query().Get("id")

	sources, err := data.LoadSources()
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	sourcesProxy := data.NewIRAProxy(sources)

	properties := []string{
		data.GetContentErrorsProperty,
		data.MatchContentErrorsProperty,
		data.ReduceErrorsProperty,
		data.CurrentElementsProperty,
		data.AddedElementsProperty,
		data.RemovedElementsProperty,
		data.SourceURLProperty}

	rdxIRA := make(kvas.IdReduxAssets)

	rdx, err := kvas.ConnectReduxAssets(data.AbsReduxDir(), properties...)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	ids := make(map[string]interface{})
	for _, p := range properties {
		for _, pid := range rdx.Keys(p) {
			ids[pid] = nil
		}
	}

	for pid := range ids {
		for _, p := range properties {
			if rdxIRA[pid] == nil {
				rdxIRA[pid] = make(map[string][]string)
			}
			rdxIRA[pid][p], _ = rdx.GetAllValues(p, pid)
			if rdxIRA[pid][p] == nil {
				rdxIRA[pid][p] = make([]string, 0)
			}
		}
	}

	sourcesProxy.Merge(rdxIRA)

	DefaultHeaders(w)

	if err := app.RenderItem(id, nil, sourcesProxy, w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

}
