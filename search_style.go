package brewerydb

// Style retrieves one or more Styles matching the given query string.
// TODO: pagination??
func (ss *SearchService) Style(query string, withDescriptions bool) ([]Style, error) {
	q := struct {
		Query            string `url:"q"`
		WithDescriptions string `url:"withDescriptions,omitempty"`
	}{Query: query}
	if withDescriptions {
		q.WithDescriptions = "Y"
	}

	req, err := ss.c.NewRequest("GET", "/search/style", &q)
	if err != nil {
		return nil, err
	}

	styleResponse := struct {
		NumberOfPages int
		CurrentPage   int
		TotalResults  int
		Data          []Style
	}{}
	err = ss.c.Do(req, &styleResponse)
	return styleResponse.Data, err
}
