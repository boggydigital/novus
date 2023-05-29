package rest

import (
	"fmt"
	"github.com/boggydigital/middleware"
	"github.com/boggydigital/nod"
	"github.com/boggydigitl/novus/data"
	"net/http"
	"os"
)

func GetAtom(w http.ResponseWriter, r *http.Request) {

	// GET /atom

	absAtomFeedPath := data.AbsAtomPath()
	if stat, err := os.Stat(absAtomFeedPath); err == nil {
		lm := stat.ModTime().UTC().Format(http.TimeFormat)
		w.Header().Set(middleware.LastModifiedHeader, lm)
		ims := r.Header.Get(middleware.IfModifiedSinceHeader)
		if middleware.IsNotModified(ims, lm) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		http.ServeFile(w, r, absAtomFeedPath)
	} else {
		_ = nod.Error(fmt.Errorf("atom feed not found"))
		http.NotFound(w, r)
	}
}
