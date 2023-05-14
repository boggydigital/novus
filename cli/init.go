package cli

import (
	"html/template"
	"io/fs"
)

var (
	tmpl *template.Template
)

func Init(templatesFS fs.FS, tmplFuncs template.FuncMap) {

	tmpl = template.Must(
		template.New("novus").
			Funcs(tmplFuncs).
			ParseFS(templatesFS, "templates/*.gohtml"))
}
