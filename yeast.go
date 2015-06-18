package brewerydb

import (
	"fmt"
	"net/http"
)

// YeastService provides access to the BreweryDB Yeast API. Use Client.Yeast.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/yeast_index
type YeastService struct {
	c *Client
}

// YeastType represents a specific type/usage of yeast.
type YeastType string

// Specific Yeast types.
const (
	YeastTypeAle       YeastType = "ale"
	YeastTypeWheat               = "wheat"
	YeastTypeLager               = "lager"
	YeastTypeWine                = "wine"
	YeastTypeChampagne           = "champagne"
)

// Yeast represents a type of yeast used in making a Beer.
type Yeast struct {
	ID                  int
	Name                string
	Description         string
	Category            string // This will always be set to "yeast"
	CategoryDisplay     string // This will always be set to "Yeast"
	YeastType           YeastType
	AttenuationMin      float64
	AttenuationMax      float64
	FermentTempMin      float64
	FermentTempMax      float64
	AlcoholToleranceMin float64
	AlcoholToleranceMax float64
	ProductID           string
	Supplier            string
	YeastFormat         string
	CreateDate          string
	UpdateDate          string
}

// YeastList represents a single "page" containing a slice of Yeasts.
type YeastList struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Yeasts        []Yeast `json:"data"`
}

// List returns all Yeasts on the given page.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/yeast_index#1
func (ys *YeastService) List(page int) (yl YeastList, err error) {
	// GET: /yeasts
	var req *http.Request
	req, err = ys.c.NewRequest("GET", "/yeasts", &Page{page})
	if err != nil {
		return
	}

	err = ys.c.Do(req, &yl)
	return
}

// Get obtains the Yeast with the given Yeast ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/yeast_index#2
func (ys *YeastService) Get(id int) (y Yeast, err error) {
	// GET: /yeast/:yeastID
	var req *http.Request
	req, err = ys.c.NewRequest("GET", fmt.Sprintf("/yeast/%d", id), nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    Yeast
		Message string
	}{}
	err = ys.c.Do(req, &resp)
	return resp.Data, err
}
