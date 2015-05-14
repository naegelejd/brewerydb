package brewerydb

// UPC retrieves one or more Beers matching the given Universal Product Code.
// TODO: pagination??
func (ss *SearchService) UPC(code string) ([]Beer, error) {
	q := struct {
		Code string `url:"code"`
	}{code}

	req, err := ss.c.NewRequest("GET", "/search/upc", &q)
	if err != nil {
		return nil, err
	}

	upcResponse := struct {
		NumberOfPages int
		CurrentPage   int
		TotalResults  int
		Data          []Beer
	}{}
	err = ss.c.Do(req, &upcResponse)
	return upcResponse.Data, err
}
