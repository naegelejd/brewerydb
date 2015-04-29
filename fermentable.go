package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// FermentableService provides access to the BreweryDB Fermentable API.
// Use Client.Fermentable.
type FermentableService struct {
	c *Client
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
	Category             string
	CategoryDisplay      string
	CreateDate           string
	UpdateDate           string
	Srm                  struct {
		ID   int
		Name string
		Hex  string
	}
	Country struct {
		IsoCode     string
		Name        string
		DisplayName string
		IsoThree    string
		NumberCode  int
		CreateDate  string
	}
	Characteristics []string
}

// FermentableList represents a "page" containing a slice of Fermentables.
type FermentableList struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Fermentables  []Fermentable `json:"data"`
}

// Fermentables returns a list of Fermentable Beer ingredients.
func (fs *FermentableService) Fermentables(page int) (fl FermentableList, err error) {
	// GET: /fermentables
	v := url.Values{}
	v.Set("p", strconv.Itoa(page))
	u := fs.c.url("/fermentables", &v)

	var resp *http.Response
	resp, err = fs.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get fermentables")
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&fl); err != nil {
		return
	}
	return
}

// Fermentable returns the Fermentable with the given Fermentable ID.
func (fs *FermentableService) Fermentable(id int) (f Fermentable, err error) {
	// GET: /fermentable/:fermentableID
	u := fs.c.url(fmt.Sprintf("/fermentable/%d", id), nil)

	var resp *http.Response
	resp, err = fs.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get fermentable")
		return
	}
	defer resp.Body.Close()

	fermentableResponse := struct {
		Status  string
		Data    Fermentable
		Message string
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&fermentableResponse); err != nil {
		return
	}
	f = fermentableResponse.Data
	return
}
