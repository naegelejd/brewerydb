package brewerydb

import (
	"fmt"
	"net/http"
)

// ListYeasts returns a slice of all Yeasts in the Beer with the given ID.
func (bs *BeerService) ListYeasts(beerID string) (al []Yeast, err error) {
	// GET: /beer/:beerId/yeasts
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+beerID+"/yeasts", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Yeast
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// AddYeast adds the Yeast with the given ID to the Beer with the given ID.
func (bs *BeerService) AddYeast(beerID string, yeastID int) error {
	// POST: /beer/:beerId/yeasts
	q := struct {
		ID int `url:"yeastId"`
	}{yeastID}
	req, err := bs.c.NewRequest("POST", "/beer/"+beerID+"/yeasts", &q)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// DeleteYeast removes the Yeast with the given ID from the Beer with the given ID.
func (bs *BeerService) DeleteYeast(beerID string, yeastID int) error {
	// DELETE: /beer/:beerId/yeast/:yeastId
	req, err := bs.c.NewRequest("DELETE", fmt.Sprintf("/beer/%s/yeast/%d", beerID, yeastID), nil)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}
