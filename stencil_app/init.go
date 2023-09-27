package stencil_app

import (
	"github.com/boggydigital/stencil"
	"github.com/boggydigitl/novus/data"
)

const (
	appTitle       = "novus"
	appAccentColor = "tomato"
)

func Init() (*stencil.AppConfiguration, error) {

	app := stencil.NewAppConfig(appTitle, appAccentColor)

	app.SetNavigation(NavItems, NavIcons, NavHrefs)
	app.SetFooter(FooterLocation, FooterRepoUrl)

	if err := app.SetCommonConfiguration(
		SourceLabels,
		nil,
		nil,
		data.Title,
		PropertyTitles,
		SectionTitles,
		nil); err != nil {
		return app, nil
	}

	if err := app.SetListConfiguration(
		ListProperties,
		SourceLabels,
		SourcePath,
		"",
		"",
		nil); err != nil {
		return app, err
	}

	if err := app.SetItemConfiguration(
		SourceProperties,
		nil,
		HiddenProperties,
		nil,
		"",
		"",
		nil); err != nil {
		return app, err
	}

	app.SetFormatterConfiguration(
		fmtLabel, fmtTitle, fmtHref, nil, fmtAction, fmtActionHref)

	return app, nil

}
