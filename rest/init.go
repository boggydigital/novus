package rest

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/novus/data"
	"github.com/boggydigital/novus/stencil_app"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/stencil"
	"html/template"
	"io/fs"
)

var (
	tmpl *template.Template
	app  *stencil.AppConfiguration
	rdx  kvas.ReadableRedux
)

func Init(templatesFS fs.FS, tmplFuncs template.FuncMap, stencilAppStyles fs.FS) error {
	tmpl = template.Must(
		template.New("").
			Funcs(tmplFuncs).
			ParseFS(templatesFS, "templates/*.gohtml"))

	stencil.InitAppTemplates(stencilAppStyles, "stencil_app/styles/css.gohtml")

	var err error
	app, err = stencil_app.Init()

	arp, err := pathways.GetAbsDir(data.Redux)
	if err != nil {
		return err
	}

	rdx, err = kvas.NewReduxReader(arp, data.AllProperties()...)
	if err != nil {
		return err
	}

	return err
}
