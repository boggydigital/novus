package data

import "github.com/boggydigital/kvas"

func NewIRAProxy(sources []*Source) *kvas.IRAProxy {
	data := kvas.IdReduxAssets{}

	for _, src := range sources {
		entry := make(map[string][]string)
		entry[URL] = []string{src.URL.String()}
		entry[Recipe] = []string{src.Recipe}
		entry[Encoding] = []string{src.Encoding}
		entry[Title] = []string{src.Title}
		entry[Category] = []string{src.Category}

		data[src.Id] = entry
	}

	return kvas.NewIRAProxy(data)
}
