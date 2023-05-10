package rest

import (
	"github.com/boggydigital/stencil"
	"github.com/boggydigitl/novus/stencil_app"
	"html/template"
	"io/fs"
)

var (
	tmpl *template.Template
	app  *stencil.AppConfiguration
)

func Init(templatesFS fs.FS, stencilAppStyles fs.FS) error {
	tmpl = template.Must(
		template.New("").
			ParseFS(templatesFS, "templates/*.gohtml"))

	stencil.InitAppTemplates(stencilAppStyles, "stencil_app/styles/css.gohtml")

	var err error
	app, err = stencil_app.Init()

	return err
}
