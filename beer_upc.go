package brewerydb

// AddUPC assigns a Universal Product Code to the Beer with the given ID.
// fluidsizeID is optional.
func (bs *BeerService) AddUPC(beerID string, code uint64, fluidsizeID *int) error {
	// POST: /beer/:beerId/upcs
	q := struct {
		Code        uint64 `url:"upcCode"`
		FluidsizeID int    `url:"fluidSizeId,omitempty"`
	}{Code: code}

	if fluidsizeID != nil {
		q.FluidsizeID = *fluidsizeID
	}

	req, err := bs.c.NewRequest("POST", "/beer/"+beerID+"/upcs", &q)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}
