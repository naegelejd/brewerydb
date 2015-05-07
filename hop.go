package brewerydb

import (
	"fmt"
	"net/http"
)

// HopService provides access to the BreweryDB Hop API. Use Client.Hop.
type HopService struct {
	c *Client
}

// HopList represents a "page" containing a slice of Hops.
type HopList struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Hops          []Hop `json:"data"`
}

// Hop contains all relevant information for a single variety of Hop.
type Hop struct {
	ID               int
	Name             string
	Description      string
	CountryOfOrigin  string
	AlphaAcidMin     float64
	AlphaAcidMax     float64
	BetaAcidMin      float64
	BetaAcidMax      float64
	HumuleneMin      float64
	HumuleneMax      float64
	CaryophylleneMin float64
	CaryophylleneMax float64
	CohumuloneMin    float64
	CohumuloneMax    float64
	MyrceneMin       float64
	MyrceneMax       float64
	FarneseneMin     float64
	FarneseneMax     float64
	IsNoble          string
	ForBittering     string
	ForFlavor        string
	ForAroma         string
	Category         string
	CategoryDisplay  string
	CreateDate       string
	UpdateDate       string
	Country          struct {
		IsoCode     string
		Name        string
		DisplayName string
		IsoThree    string
		NumberCode  int
		CreateDate  string
	}
}

// List returns all Hops on the given page.
func (hs *HopService) List(page int) (hl HopList, err error) {
	var req *http.Request
	req, err = hs.c.NewRequest("GET", "/hops", Page{page})
	if err != nil {
		return
	}

	err = hs.c.Do(req, &hl)
	return
}

// Get queries for a single Hop with the given Hop ID.
func (hs *HopService) Get(id int) (hop Hop, err error) {
	// GET: /hop/:hopId
	var req *http.Request
	req, err = hs.c.NewRequest("GET", fmt.Sprintf("/hop/%d", id), nil)
	if err != nil {
		return
	}

	hopResp := struct {
		Message string
		Data    Hop
		Status  string
	}{}

	err = hs.c.Do(req, &hopResp)
	return hopResp.Data, err
}
