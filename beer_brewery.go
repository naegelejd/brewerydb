package brewerydb

// Breweries queries for all Breweries associated with the Beer having the given ID.
func (bs *BeerService) Breweries(id string) ([]Brewery, error) {
	// GET: /beer/:beerId/breweries
	req, err := bs.c.NewRequest("GET", "/beer/"+id+"/breweries", nil)
	if err != nil {
		return nil, err
	}

	resp := struct {
		Status  string
		Data    []Brewery
		Message string
	}{}
	if err := bs.c.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// BeerBreweryRequest allows for specifying locations for a given Brewery, e.g.
// if only adding/removing a specific Brewery location from a Beer.
type BeerBreweryRequest struct {
	LocationID string `url:"locationId,omitempty"`
}

// AddBrewery adds the Brewery with the given Brewery ID to the Beer
// with the given Beer ID.
func (bs *BeerService) AddBrewery(beerID, breweryID string, q *BeerBreweryRequest) error {
	// POST: /beer/:beerId/brewery/:breweryId
	params := struct {
		ID         string `url:"breweryId"`
		LocationID string `url:"locationId,omitempty"`
	}{ID: breweryID}

	if q != nil {
		params.LocationID = q.LocationID
	}

	req, err := bs.c.NewRequest("POST", "/beer/"+beerID+"/brewery/"+breweryID, &params)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// DeleteBrewery removes the Brewery with the given Brewery ID from the Beer
// with the given Beer ID.
func (bs *BeerService) DeleteBrewery(beerID, breweryID string, q *BeerBreweryRequest) error {
	// DELETE: /beer/:beerId/brewery/:breweryId
	req, err := bs.c.NewRequest("DELETE", "/beer/"+beerID+"/brewery/"+breweryID, q)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}
