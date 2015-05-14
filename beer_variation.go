package brewerydb

import "net/http"

// ListVariations returns a slice of all Beers that are variations of the
// Beer with the given ID.
func (bs *BeerService) ListVariations(beerID string) (bl []Beer, err error) {
	// GET: /beer/:beerId/variations
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+beerID+"/variations", nil)
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
