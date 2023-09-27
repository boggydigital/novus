package rest

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigitl/novus/data"
	"github.com/boggydigitl/novus/stencil_app"
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

	sourcesProxy := data.NewIRAProxy(sources)

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

	if err := app.RenderGroup(
		stencil_app.NavSources,
		categories,
		sourcesByCategory,
		categoriesTitles,
		categoryTotals,
		"",
		nil,
		sourcesProxy,
		w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}
