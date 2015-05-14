package brewerydb

import (
	"fmt"
	"net/http"
)

// ListHops returns a slice of all Hops in the Beer with the given ID.
func (bs *BeerService) ListHops(beerID string) (al []Hop, err error) {
	// GET: /beer/:beerId/hops
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+beerID+"/hops", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Hop
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// AddHop adds the Hop with the given ID to the Beer with the given ID.
func (bs *BeerService) AddHop(beerID string, hopID int) error {
	// POST: /beer/:beerId/hops
	q := struct {
		ID int `url:"hopId"`
	}{hopID}
	req, err := bs.c.NewRequest("POST", "/beer/"+beerID+"/hops", &q)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// DeleteHop removes the Hop with the given ID from the Beer with the given ID.
func (bs *BeerService) DeleteHop(beerID string, hopID int) error {
	// DELETE: /beer/:beerId/hop/:hopId
	req, err := bs.c.NewRequest("DELETE", fmt.Sprintf("/beer/%s/hop/%d", beerID, hopID), nil)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}
