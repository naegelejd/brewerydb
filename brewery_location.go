package brewerydb

import "net/http"

// ListLocations returns a slice of all locations for the Brewery with the given ID.
func (bs *BreweryService) ListLocations(breweryID string) (ll []Location, err error) {
	// GET: /brewery/:breweryId/locations
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/brewery/"+breweryID+"/locations", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Location
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// AddLocation adds a new location for the Brewery with the given ID.
func (bs *BreweryService) AddLocation(breweryID string, loc *Location) error {
	// POST: /brewery/:breweryId/locations
	req, err := bs.c.NewRequest("POST", "/brewery/"+breweryID+"/locations", loc)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}
