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

	srcIpv := data.SourcesIdPropertyValues(sources)

	properties := []string{
		data.GetContentErrorsProperty,
		data.MatchContentErrorsProperty,
		data.ReduceErrorsProperty,
		data.CurrentElementsProperty,
		data.AddedElementsProperty,
		data.RemovedElementsProperty,
		data.SourceURLProperty}

	rdx, err := kvas.NewReduxReader(data.AbsReduxDir(), properties...)
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
			if srcIpv[pid] == nil {
				srcIpv[pid] = make(map[string][]string)
			}
			srcIpv[pid][p], _ = rdx.GetAllValues(p, pid)
			if srcIpv[pid][p] == nil {
				srcIpv[pid][p] = make([]string, 0)
			}
		}
	}

	rdxProxy := kvas.ReduxProxy(srcIpv)

	DefaultHeaders(w)

	if err := app.RenderItem(id, nil, rdxProxy, w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

}
