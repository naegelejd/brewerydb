package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Breweries queries for all Breweries associated with the Beer having the given ID.
// GET: /beer/:beerId/breweries
func (s *BeerService) Breweries(id string) ([]*Brewery, error) {
	u := s.c.url("/beer/"+id+"/breweries", nil)

	resp, err := s.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("breweries not found")
	}
	defer resp.Body.Close()

	breweriesResp := struct {
		Status  string
		Data    []*Brewery
		Message string
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&breweriesResp); err != nil {
		return nil, err
	}

	return breweriesResp.Data, nil
}

// BeerBreweryRequest allows for specifying locations for a given Brewery, e.g.
// if only adding/removing a specific Brewery location from a Beer.
type BeerBreweryRequest struct {
	LocationID string `json:"locationId"`
}

// AddBrewery adds the Brewery with the given Brewery ID to the Beer
// with the given Beer ID.
// WRONG in documentation: POST: /beer/:beerId/breweries
// POST: /beer/:beerId/brewery/:breweryId
func (s *BeerService) AddBrewery(beerID, breweryID string, req *BeerBreweryRequest) error {
	u := s.c.url("/beer/"+beerID+"/brewery/"+breweryID, nil)

	resp, err := s.c.PostForm(u, encode(req))
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to add brewery")
	}
	defer resp.Body.Close()

	return nil
}

// DeleteBrewery removes the Brewery with the given Brewery ID from the Beer
// with the given Beer ID.
// DELETE: /beer/:beerId/brewery/:breweryId
func (s *BeerService) DeleteBrewery(beerID, breweryID string, req *BeerBreweryRequest) error {
	v := encode(req)
	u := s.c.url("/beer/"+beerID+"/brewery/"+breweryID, &v)
	q, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	resp, err := s.c.Do(q)
	if resp != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to delete brewery")
	}
	defer resp.Body.Close()

	return nil
}
