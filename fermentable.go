package brewerydb

import (
	"fmt"
	"net/http"
)

// FermentableService provides access to the BreweryDB Fermentable API.
// Use Client.Fermentable.
type FermentableService struct {
	c *Client
}

// Characteristic represents a descriptive characteristic of a Fermentable.
type Characteristic struct {
	ID          int
	Name        string
	Description string
	CreateDate  string
	UpdateDate  string
}

// Fermentable represents a Fermentable Beer ingredient.
type Fermentable struct {
	ID                   int
	Name                 string
	Description          string
	CountryOfOrigin      string
	SrmID                int
	SrmPrecise           float64
	MoistureContent      float64
	CoarseFineDifference float64
	DiastaticPower       float64
	DryYield             float64
	Potential            float64
	Protein              float64
	SolubleNitrogenRatio float64
	MaxInBatch           float64
	RequiresMashing      string // Y or N
	Category             string // This will always be set to "malt"
	CategoryDisplay      string // This will always be set to "Malts, Grains, & Fermentables"
	CreateDate           string
	UpdateDate           string
	SRM                  SRM
	Country              struct {
		IsoCode     string
		Name        string
		DisplayName string
		IsoThree    string
		NumberCode  int
		CreateDate  string
	}
	Characteristics []Characteristic
}

// FermentableList represents a "page" containing a slice of Fermentables.
type FermentableList struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Fermentables  []Fermentable `json:"data"`
}

// List returns a list of Fermentable Beer ingredients.
func (fs *FermentableService) List(page int) (fl FermentableList, err error) {
	// GET: /fermentables
	var req *http.Request
	req, err = fs.c.NewRequest("GET", "/fermentables", &Page{page})
	if err != nil {
		return
	}

	err = fs.c.Do(req, &fl)
	return
}

// Get returns the Fermentable with the given Fermentable ID.
func (fs *FermentableService) Get(id int) (f Fermentable, err error) {
	// GET: /fermentable/:fermentableID
	var req *http.Request
	req, err = fs.c.NewRequest("GET", fmt.Sprintf("/fermentable/%d", id), nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    Fermentable
		Message string
	}{}
	err = fs.c.Do(req, &resp)
	return resp.Data, err
}
