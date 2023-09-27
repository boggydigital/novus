package stencil_app

import "github.com/boggydigitl/novus/data"

var SourceProperties = []string{
	data.Title,
	data.URL,
	data.Recipe,
	data.Encoding,
	data.Category,
	data.MatchContentErrorsProperty,
	data.ReduceErrorsProperty,
	data.CurrentElementsProperty,
	data.AddedElementsProperty,
	data.RemovedElementsProperty,
}

var HiddenProperties = []string{
	data.Title,
	data.Category,
}
