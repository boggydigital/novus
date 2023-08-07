package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigitl/novus/data"
	"github.com/boggydigitl/novus/stencil_app"
	"html/template"
	"net/http"
	"sort"
	"strings"
)

type Sources struct {
	Ids             []string
	Titles          map[string]string
	Categories      map[string]string
	Hosts           map[string]string
	URLs            map[string]string
	CurrentElements map[string][]template.HTML
}

func GetSources(w http.ResponseWriter, r *http.Request) {

	// GET /sources

	rdx, err := kvas.ConnectReduxAssets(data.AbsReduxDir(),
		data.CurrentElementsProperty,
		data.SourceURLProperty)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	sb := &strings.Builder{}

	sources, err := data.LoadSources()
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	ids := data.SourcesIds(sources...)
	sort.Strings(ids)

	svm := &Sources{
		Ids:             ids,
		Titles:          make(map[string]string),
		Categories:      make(map[string]string),
		Hosts:           make(map[string]string),
		URLs:            make(map[string]string),
		CurrentElements: make(map[string][]template.HTML),
	}

	for _, id := range ids {
		if su, ok := rdx.GetFirstVal(data.SourceURLProperty, id); ok {
			svm.Hosts[id] = Host(su)
		}

		src := data.SourceById(id, sources...)
		svm.URLs[id] = src.URL.String()
		svm.Titles[id] = src.Title
		svm.Categories[id] = src.Category

		if currentElements, ok := rdx.GetAllValues(data.CurrentElementsProperty, src.Id); ok {
			for _, ce := range currentElements {
				svm.CurrentElements[id] = append(svm.CurrentElements[id], template.HTML(AbsHref(ce, svm.Hosts[id])))
			}
		}
	}

	if err := tmpl.ExecuteTemplate(sb, "sources", svm); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	DefaultHeaders(w)

	if err := app.RenderPage("", stencil_app.NavSources, sb.String(), w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}
