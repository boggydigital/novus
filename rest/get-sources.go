package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigitl/novus/data"
	"github.com/boggydigitl/novus/stencil_app"
	"net/http"
	"sort"
	"strings"
)

func GetSources(w http.ResponseWriter, r *http.Request) {

	// GET /sources

	rdx, err := kvas.ConnectReduxAssets(data.AbsReduxDir(), nil,
		data.CurrentElementsProperty,
		data.SourceURLProperty)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	sb := strings.Builder{}

	sources, err := data.LoadSources()
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	ids := data.SourcesIds(sources...)
	sort.Strings(ids)

	sb.WriteString("<section>")
	for _, id := range ids {
		host := ""
		if su, ok := rdx.GetFirstVal(data.SourceURLProperty, id); ok {
			host = Host(su)
		}

		src := data.SourceById(id, sources...)

		sb.WriteString("<details>")
		sb.WriteString("<summary>" + src.Id + "</summary>")
		sb.WriteString("<a href='" + src.URL.String() + "'>Source Link</a>")

		if currentElements, ok := rdx.GetAllUnchangedValues(data.CurrentElementsProperty, src.Id); ok {
			sb.WriteString("<ul>")
			for _, ce := range currentElements {
				sb.WriteString("<li>" + AbsHref(ce, host) + "</li>")
			}
			sb.WriteString("</ul>")
		}

		sb.WriteString("</details>")
	}
	sb.WriteString("</section>")

	DefaultHeaders(w)

	if err := app.RenderPage("", stencil_app.NavSources, sb.String(), w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}
