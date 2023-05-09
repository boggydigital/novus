package data

import (
	"bytes"
	_ "embed"
	"github.com/boggydigital/wits"
)

//go:embed "builtin-recipes.txt"
var recipesDefinitions []byte

func LoadRecipes() (wits.SectionKeyValue, error) {
	return wits.ReadSectionKeyValue(bytes.NewReader(recipesDefinitions))
}
