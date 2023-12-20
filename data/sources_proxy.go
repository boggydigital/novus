package data

func SourcesIdPropertyValues(sources []*Source) map[string]map[string][]string {
	data := make(map[string]map[string][]string)

	for _, src := range sources {
		entry := make(map[string][]string)
		entry[URL] = []string{src.URL.String()}
		entry[Recipe] = []string{src.Recipe}
		entry[Encoding] = []string{src.Encoding}
		entry[Title] = []string{src.Title}
		entry[Category] = []string{src.Category}

		data[src.Id] = entry
	}

	return data
}
