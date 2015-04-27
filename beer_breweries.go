package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GET: /beer/:beerId/breweries
func (c *Client) BeerBreweries(id string) ([]*Brewery, error) {
	u := c.url("/beer/"+id+"/breweries", nil)

	resp, err := c.Get(u)
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

type BeerBreweryRequest struct {
	LocationID string `json:"locationId"`
}

// WRONG in documentation: POST: /beer/:beerId/breweries
// POST: /beer/:beerId/brewery/:breweryId
// TODO: rename to AddBreweryToBeer?
func (c *Client) AddBeerBrewery(beerID, breweryID string, req *BeerBreweryRequest) error {
	u := c.url("/beer/"+beerID+"/brewery/"+breweryID, nil)

	resp, err := c.PostForm(u, encode(req))
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to add brewery")
	}
	defer resp.Body.Close()

	return nil
}

// DELETE: /beer/:beerId/brewery/:breweryId
// TODO: rename to DeleteBreweryFromBeer
func (c *Client) DeleteBeerBreweries(beerID, breweryID string, req *BeerBreweryRequest) error {
	v := encode(req)
	u := c.url("/beer/"+beerID+"/brewery/"+breweryID, &v)
	q, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(q)
	if resp != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to delete brewery")
	}
	defer resp.Body.Close()

	return nil
}
