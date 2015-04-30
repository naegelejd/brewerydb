package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GlassService provides access to the BreweryDB Glassware API.
// Use Client.Glass.
type GlassService struct {
	c *Client
}

// Glass represents a Glass assigned to a UPC code.
type Glass struct {
	ID          int
	Name        string
	Description string
	CreateDate  string
	UpdateDate  string
}

// List returns a list of Glasses.
func (gs *GlassService) List() (gl []Glass, err error) {
	// GET: /glassware
	u := gs.c.url("/glassware", nil)

	var resp *http.Response
	resp, err = gs.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get glassware")
		return
	}
	defer resp.Body.Close()

	glasswareResponse := struct {
		Status  string
		Data    []Glass
		Message string
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&glasswareResponse); err != nil {
		return
	}
	gl = glasswareResponse.Data
	return
}

// Get returns the Glass with the given Glass ID.
func (gs *GlassService) Get(id int) (g Glass, err error) {
	// GET: /glass/:glassId
	u := gs.c.url(fmt.Sprintf("/glass/%d", id), nil)

	var resp *http.Response
	resp, err = gs.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get glass")
		return
	}
	defer resp.Body.Close()

	glassResponse := struct {
		Status  string
		Data    Glass
		Message string
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&glassResponse); err != nil {
		return
	}
	g = glassResponse.Data
	return
}
