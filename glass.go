package brewerydb

import (
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
	var req *http.Request
	req, err = gs.c.NewRequest("GET", "/glassware", nil)
	if err != nil {
		return
	}

	glasswareResponse := struct {
		Status  string
		Data    []Glass
		Message string
	}{}
	if err = gs.c.Do(req, &glasswareResponse); err != nil {
		return
	}

	return glasswareResponse.Data, nil
}

// Get returns the Glass with the given Glass ID.
func (gs *GlassService) Get(id int) (g Glass, err error) {
	// GET: /glass/:glassId
	var req *http.Request
	req, err = gs.c.NewRequest("GET", fmt.Sprintf("/glass/%d", id), nil)
	if err != nil {
		return
	}

	glassResponse := struct {
		Status  string
		Data    Glass
		Message string
	}{}

	if err = gs.c.Do(req, &glassResponse); err != nil {
		return
	}
	return glassResponse.Data, nil
}
