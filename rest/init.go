package rest

import (
	"github.com/boggydigital/novus/stencil_app"
	"github.com/boggydigital/stencil"
	"html/template"
	"io/fs"
)

var (
	tmpl *template.Template
	app  *stencil.AppConfiguration
)

func Init(templatesFS fs.FS, tmplFuncs template.FuncMap, stencilAppStyles fs.FS) error {
	tmpl = template.Must(
		template.New("").
			Funcs(tmplFuncs).
			ParseFS(templatesFS, "templates/*.gohtml"))

	stencil.InitAppTemplates(stencilAppStyles, "stencil_app/styles/css.gohtml")

	var err error
	app, err = stencil_app.Init()

	return err
}
