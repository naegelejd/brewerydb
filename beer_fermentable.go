package brewerydb

import (
	"fmt"
	"net/http"
)

// ListFermentables returns a slice of all Fermentables in the Beer with the given ID.
func (bs *BeerService) ListFermentables(beerID string) (fl []Fermentable, err error) {
	// GET: /beer/:beerId/fermentables
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+beerID+"/fermentables", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Fermentable
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// AddFermentable adds the Fermentable with the given ID to the Beer with the given ID.
func (bs *BeerService) AddFermentable(beerID string, fermentableID int) error {
	// POST: /beer/:beerId/fermentables
	q := struct {
		ID int `url:"fermentableId"`
	}{fermentableID}
	req, err := bs.c.NewRequest("POST", "/beer/"+beerID+"/fermentables", &q)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}

// DeleteFermentable removes the Fermentable with the given ID from the Beer with the given ID.
func (bs *BeerService) DeleteFermentable(beerID string, fermentableID int) error {
	// DELETE: /beer/:beerId/fermentable/:fermentableId
	req, err := bs.c.NewRequest("DELETE", fmt.Sprintf("/beer/%s/fermentable/%d", beerID, fermentableID), nil)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}
