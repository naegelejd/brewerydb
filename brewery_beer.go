package brewerydb

import "net/http"

// BreweryBeersRequest contains options for querying for all Beers from a Brewery.
type BreweryBeersRequest struct {
	WithBreweries      string `url:"withBreweries"`      // Y/N
	WithSocialAccounts string `url:"withSocialAccounts"` // Y/N
	WithIngredients    string `url:"withIngredients"`    // Y/N
}

// ListBeers returns a slice of all Beers offered by the Brewery with the given ID.
func (bs *BreweryService) ListBeers(breweryID string, q *BreweryBeersRequest) (bl []Beer, err error) {
	// GET: /brewery/:breweryId/beers
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/brewery/"+breweryID+"/beers", q)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Beer
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}
