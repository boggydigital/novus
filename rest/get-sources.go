package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/novus/data"
	"github.com/boggydigital/novus/stencil_app"
	"golang.org/x/exp/maps"
	"net/http"
	"sort"
)

func GetSources(w http.ResponseWriter, r *http.Request) {

	// GET /sources

	sources, err := data.LoadSources()
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	srcIpv := data.SourcesIdPropertyValues(sources)

	DefaultHeaders(w)

	sourcesByCategory := make(map[string][]string)
	for _, src := range sources {
		sourcesByCategory[src.Category] = append(sourcesByCategory[src.Category], src.Id)
	}

	categories := maps.Keys(sourcesByCategory)
	sort.Strings(categories)
	for category, _ := range sourcesByCategory {
		sort.Strings(sourcesByCategory[category])
	}

	// sorting categories and sources
	categoriesTitles := make(map[string]string)
	categoryTotals := make(map[string]int)
	for _, c := range categories {
		categoriesTitles[c] = c
		categoryTotals[c] = len(sourcesByCategory[c])
	}

	rdx := kvas.ReduxProxy(srcIpv)

	if err := app.RenderGroup(
		stencil_app.NavSources,
		categories,
		sourcesByCategory,
		categoriesTitles,
		categoryTotals,
		"",
		nil,
		rdx,
		w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}
