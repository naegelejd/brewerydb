package brewerydb

import (
	"fmt"
	"net/http"
)

// ListAdjuncts returns a slice of all Adjuncts in the Beer with the given ID.
func (bs *BeerService) ListAdjuncts(beerID string) (al []Adjunct, err error) {
	// GET: /beer/:beerId/adjuncts
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+beerID+"/adjuncts", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Adjunct
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// AddAdjunct adds the Adjunct with the given ID to the Beer with the given ID.
func (bs *BeerService) AddAdjunct(beerID string, adjunctID int) error {
	// POST: /beer/:beerId/adjuncts
	q := struct {
		ID int `url:"adjunctId"`
	}{adjunctID}
	req, err := bs.c.NewRequest("POST", "/beer/"+beerID+"/adjuncts", &q)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// DeleteAdjunct removes the Adjunct with the given ID from the Beer with the given ID.
func (bs *BeerService) DeleteAdjunct(beerID string, adjunctID int) error {
	// DELETE: /beer/:beerId/adjunct/:adjunctId
	req, err := bs.c.NewRequest("DELETE", fmt.Sprintf("/beer/%s/adjunct/%d", beerID, adjunctID), nil)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}
