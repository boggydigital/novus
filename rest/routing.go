package rest

import (
	"github.com/boggydigital/nod"
	"net/http"
)

var (
	Log = nod.RequestLog
)

func HandleFuncs() {
	patternHandlers := map[string]http.Handler{
		"GET /atom":    Log(http.HandlerFunc(GetAtom)),
		"GET /sources": Log(http.HandlerFunc(GetSources)),
		"GET /source":  Log(http.HandlerFunc(GetSource)),
		"GET /":        http.RedirectHandler("/sources", http.StatusPermanentRedirect),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
