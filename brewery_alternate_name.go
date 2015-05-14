package brewerydb

import (
	"fmt"
	"net/http"
)

// AlternateName represents an alternate name for a Brewery.
// TODO: the actual response object contains the entire Brewery object as well.
//	see: http://www.brewerydb.com/developers/docs-endpoint/brewery_alternatename
type AlternateName struct {
	ID         int
	Name       string
	BreweryID  string
	CreateDate string
	UpdateDate string
}

// ListAlternateNames returns a slice of all the AlternateNames for the Brewery with the given ID.
func (bs *BreweryService) ListAlternateNames(breweryID string) (al []AlternateName, err error) {
	// GET: /brewery/:breweryId/alternatenames
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/brewery/"+breweryID+"/alternatenames", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []AlternateName
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// AddAlternateName adds an alternate name to the Brewery with the given ID.
func (bs *BreweryService) AddAlternateName(breweryID, name string) error {
	// POST: /brewery/:breweryId/alternatenames
	q := struct {
		Name string `url:"name"`
	}{name}
	req, err := bs.c.NewRequest("POST", "/brewery/"+breweryID+"/alternatenames", &q)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}

// DeleteAlternateName removes the AlternateName with the given ID from the Brewery with the given ID.
func (bs *BreweryService) DeleteAlternateName(breweryID string, alternateNameID int) error {
	// DELETE: /brewery/:breweryId/alternatename/:alternatenameId
	req, err := bs.c.NewRequest("DELETE", fmt.Sprintf("/brewery/%s/alternatenames/%d", breweryID, alternateNameID), nil)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}
