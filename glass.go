package brewerydb

import (
	"fmt"
	"net/http"
)

// GlassService provides access to the BreweryDB Glassware API.
// Use Client.Glass.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/glass_index
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
//
// See: http://www.brewerydb.com/developers/docs-endpoint/glass_index#1
func (gs *GlassService) List() (gl []Glass, err error) {
	// GET: /glassware
	var req *http.Request
	req, err = gs.c.NewRequest("GET", "/glassware", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Glass
		Message string
	}{}
	err = gs.c.Do(req, &resp)
	return resp.Data, err
}

// Get returns the Glass with the given Glass ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/glass_index#2
func (gs *GlassService) Get(id int) (g Glass, err error) {
	// GET: /glass/:glassId
	var req *http.Request
	req, err = gs.c.NewRequest("GET", fmt.Sprintf("/glass/%d", id), nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    Glass
		Message string
	}{}
	err = gs.c.Do(req, &resp)
	return resp.Data, err
}
