package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigitl/novus/data"
	"io"
	"net/http"
	"sort"
	"strings"
)

func GetSources(w http.ResponseWriter, r *http.Request) {

	// GET /sources

	rdx, err := kvas.ConnectRedux(data.AbsReduxDir(), data.CurrentElementsProperty)
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

	for _, id := range ids {

		src := data.SourceById(id, sources...)

		sb.WriteString("<details>")
		sb.WriteString("<summary>" + src.Id + "</summary>")
		sb.WriteString("<a href='" + src.URL.String() + "'>URL</a>")

		if currentElements, ok := rdx.GetAllValues(src.Id); ok {
			sb.WriteString("<ul>")
			for _, ce := range currentElements {
				sb.WriteString("<li>" + ce + "</li>")
			}
			sb.WriteString("</ul>")
		}

		sb.WriteString("</details>")
	}

	w.Header().Set("Content-Type", "text/html")

	if _, err := io.Copy(w, strings.NewReader(sb.String())); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}
