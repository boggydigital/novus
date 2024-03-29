package stencil_app

import "github.com/boggydigital/novus/data"

var PropertyTitles = map[string]string{
	data.Id:                         "Id",
	data.URL:                        "Source",
	data.Recipe:                     "Recipe",
	data.Encoding:                   "Encoding",
	data.Title:                      "Title",
	data.Category:                   "Category",
	data.CurrentElementsProperty:    "Current Elements",
	data.AddedElementsProperty:      "Added Elements",
	data.RemovedElementsProperty:    "Removed Elements",
	data.GetContentErrorsProperty:   "Get Content Errors",
	data.ReduceErrorsProperty:       "Reduce Errors",
	data.MatchContentErrorsProperty: "Match Content Errors",
}
