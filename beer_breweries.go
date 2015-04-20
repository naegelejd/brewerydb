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

// POST: /beer/:beerId/breweries
func (c *Client) AddBeerBreweries(id string) {

}

// DELETE: /beer/:beerId/brewery/:breweryId
func (c *Client) DeleteBeerBreweries(id string) {

}
