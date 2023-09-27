package stencil_app

import "github.com/boggydigitl/novus/data"

var SourceProperties = []string{
	data.Title,
	data.URL,
	data.Recipe,
	data.Encoding,
	data.Category,
	data.CurrentElementsProperty,
}

var HiddenProperties = []string{
	data.Title,
	data.Category,
}
