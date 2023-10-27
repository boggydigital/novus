package stencil_app

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/novus/data"
)

func fmtLabel(_, property, link string, _ kvas.ReduxAssets) string {
	label := link
	switch property {
	default:
		// do nothing
	}
	return label
}

func fmtTitle(id, property, link string, rxa kvas.ReduxAssets) string {
	title := link

	switch property {
	case data.URL:
		return "URL"
	default:
		// do nothing
	}

	return title
}

func fmtHref(_, property, link string, _ kvas.ReduxAssets) string {
	switch property {
	case data.URL:
		return link
	default:
		// do nothing
	}
	return ""
}

func fmtAction(id, property, link string, rxa kvas.ReduxAssets) string {
	return ""
}

func fmtActionHref(id, property, link string, _ kvas.ReduxAssets) string {
	return ""
}
