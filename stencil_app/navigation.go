package stencil_app

import "github.com/boggydigital/stencil"

const (
	NavSources = "Sources"
	NavFeed    = "Feed"
)

var NavItems = []string{NavSources, NavFeed}

var NavIcons = map[string]string{
	NavSources: stencil.IconStack,
	NavFeed:    stencil.IconSparkle,
}

var NavHrefs = map[string]string{
	NavSources: SourcesPath,
	NavFeed:    FeedPath,
}
