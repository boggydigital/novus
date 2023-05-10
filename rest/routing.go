package rest

import (
	"github.com/boggydigital/middleware"
	"github.com/boggydigital/nod"
	"net/http"
)

func HandleFuncs() {
	patternHandlers := map[string]http.Handler{
		"/atom":    middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetAtom))),
		"/sources": middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetSources))),
		"/":        http.RedirectHandler("/sources", http.StatusPermanentRedirect),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
