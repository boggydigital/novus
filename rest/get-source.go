package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigitl/novus/data"
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

	for _, p := range properties {
		for _, id := range rdx.Keys(p) {
			if rdxIRA[id] == nil {
				rdxIRA[id] = make(map[string][]string)
			}
			rdxIRA[id][p], _ = rdx.GetAllValues(p, id)
		}
	}

	sourcesProxy.Merge(rdxIRA)

	DefaultHeaders(w)

	if err := app.RenderItem(id, nil, sourcesProxy, w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

}
