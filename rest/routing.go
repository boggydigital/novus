package rest

import (
	"github.com/boggydigital/middleware"
	"github.com/boggydigital/nod"
	"net/http"
)

var (
	BrGzip  = middleware.BrGzip
	GetOnly = middleware.GetMethodOnly
	Log     = nod.RequestLog
)

func HandleFuncs() {
	patternHandlers := map[string]http.Handler{
		"/atom":    BrGzip(GetOnly(Log(http.HandlerFunc(GetAtom)))),
		"/sources": BrGzip(GetOnly(Log(http.HandlerFunc(GetSources)))),
		"/source":  BrGzip(GetOnly(Log(http.HandlerFunc(GetSource)))),
		"/":        http.RedirectHandler("/sources", http.StatusPermanentRedirect),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
