package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GET: /beer/:beerId/adjuncts
// POST: /beer/:beerId/adjuncts
// DELETE: /beer/:beerId/adjunct/:adjunctId

// GET: /beer/:beerId/events

// GET: /beer/:beerId/fermentables
// POST: /beer/:beerId/fermentables
// DELETE: /beer/:beerId/fermentable/:fermentableId

// GET: /beer/:beerId/hops
// POST: /beer/:beerId/hops
// DELETE: /beer/:beerId/hop/:hopId

// GET: /beer/:beerId/ingredients

// GET: /beer/:beerId/socialaccounts
// GET: /beer/:beerId/socialaccount/:socialaccountId
// POST: /beer/:beerId/socialaccounts
// DELETE: /beer/:beerId/socialaccount/:socialaccountId
// DELETE: /beer/:beerId/socialaccount/:socialaccountId

// POST: /beer/:beerId/upcs

// GET: /beer/:beerId/variations

// GET: /beer/:beerId/yeasts
// POST: /beer/:beerId/yeasts
// DELETE: /beer/:beerId/yeast/:yeastId

// BeerService provides access to the BreweryDB Beer API. Use Client.Beer.
type BeerService struct {
	c *Client
}

// BeerList represents a lazy list of Beers. Create a new one with
// NewBeerList. Iterate over a BeerList using First() and Next().
type BeerList struct {
	service *BeerService
	req     *BeerListRequest
	resp    *beerListResponse
	curBeer int
}

// BeerOrder represents the ordering of a list of Beers.
type BeerOrder string

// BeerList ordering options.
const (
	BeerOrderName        BeerOrder = "name"
	BeerOrderDescription           = "description"
	BeerOrderAbv                   = "abv"
	BeerOrderIbu                   = "ibu"
	BeerOrderGlasswareID           = "glasswareId"
	BeerOrderSrmID                 = "smrID"
	BeerOrderAvailableID           = "availableId"
	BeerOrderStyleID               = "styleId"
	BeerOrderIsOrganic             = "isOrganic"
	BeerOrderStatus                = "status"
	BeerOrderCreateDate            = "createDate"
	BeerOrderUpdateDate            = "updateDate"
	BeerOrderRandom                = "random"
)

// BeerTemperature represents the approximate temperature for a Beer.
type BeerTemperature string

// Beer temperatures.
const (
	TemperatureCellar   BeerTemperature = "cellar"
	TemperatureVeryCold                 = "very_cold"
	TemperatureCool                     = "cool"
	TemperatureCold                     = "cold"
	TemperatureWarm                     = "warm"
	TemperatureHot                      = "hot"
)

// BeerListRequest contains all the required and optional fields
// used for querying for a list of Beers.
type BeerListRequest struct {
	IDs                string    `json:"ids"` // IDs of the beers to return, comma separated. Max 10.
	Name               string    `json:"name"`
	ABV                string    `json:"abv"`
	IBU                string    `json:"ibu"`
	GlasswareID        int       `json:"glasswareId"`
	SrmID              int       `json:"srmId"`
	AvailableID        int       `json:"availableId"`
	StyleID            int       `json:"styleId"`
	IsOrganic          string    `json:"isOrganic"` // Y/N
	HasLabels          string    `json:"hasLabels"` // Y/N
	Year               int       `json:"year"`      // YYYY
	Since              string    `json:"since"`     // UNIX timestamp format. Max 30 days
	Status             string    `json:"status"`
	Order              BeerOrder `json:"order"`
	Sort               ListSort  `json:"sort"`
	RandomCount        string    `json:"randomCount"`        // how many random beers to return. Max 10
	WithBreweries      string    `json:"withBreweries"`      // Y/N
	WithSocialAccounts string    `json:"withSocialAccounts"` // Premium. Y/N
	WithIngredients    string    `json:"withIngredients"`    // Premium. Y/N
}

type beerListResponse struct {
	Status        string
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Beers         []Beer `json:"data"`
}

// Availability contains information on a Beer's availability.
type Availability struct {
	ID          int
	Name        string
	Description string
}

// SRM represents a Standard Reference Method.
type SRM struct {
	ID   int
	Hex  string
	Name string
}

// Beer contains all relevant information for a single Beer.
type Beer struct {
	ID              string
	Name            string
	Description     string
	FoodPairings    string
	OriginalGravity string
	ABV             string
	IBU             string
	GlasswareID     int
	Glass           Glass
	StyleID         int
	Style           Style
	IsOrganic       string
	Labels          struct {
		Medium string
		Large  string
		Icon   string
	}
	ServingTemperature        BeerTemperature
	ServingTemperatureDisplay string
	Status                    string
	StatusDisplay             string
	AvailableID               int
	Available                 Availability
	BeerVariationID           string
	BeerVariation             struct {
		// TODO: instance of a Beer??
	}
	SrmID int
	SRM   SRM
	Year  int
}

// NewBeerList returns a new BeerList that will use the given BeerListRequest
// to query for a list of Beers.
func (bs *BeerService) NewBeerList(req *BeerListRequest) *BeerList {
	// GET: /beers
	return &BeerList{service: bs, req: req}
}

// getPage obtains the "next" page from the BreweryDB API
func (bl *BeerList) getPage(pageNum int) error {
	var v url.Values
	if bl.req != nil {
		v = encode(bl.req)
	} else {
		v = url.Values{}
	}
	v.Set("p", fmt.Sprintf("%d", pageNum))

	u := bl.service.c.url("/beers", &v)

	resp, err := bl.service.c.Get(u)
	if err != nil {
		return err
	} else if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("beers not found")
	}
	defer resp.Body.Close()

	beerListResp := &beerListResponse{}
	if err := json.NewDecoder(resp.Body).Decode(beerListResp); err != nil {
		// if e, ok := err.(*json.UnmarshalTypeError); ok == true {
		// 	fmt.Printf("(JSON error) Value: %s, Type: %s", e.Value, e.Type.Kind())
		// }
		return err
	}

	if len(beerListResp.Beers) <= 0 {
		return fmt.Errorf("no beers found on page %d", pageNum)
	}

	bl.resp = beerListResp
	bl.curBeer = 0

	return nil
}

// First returns the first Beer in the BeerList.
func (bl *BeerList) First() (*Beer, error) {
	// If we already have page 1 cached, just return the first Beer
	if bl.resp != nil && bl.resp.CurrentPage == 1 {
		bl.curBeer = 0
		return &bl.resp.Beers[0], nil
	}

	if err := bl.getPage(1); err != nil {
		return nil, err
	}

	return &bl.resp.Beers[0], nil
}

// Next returns the next Beer in the BeerList on each successive call, or nil
// if there are no more Beers.
func (bl *BeerList) Next() (*Beer, error) {
	bl.curBeer++
	// if we're still on the same page just return beer
	if bl.curBeer < len(bl.resp.Beers) {
		return &bl.resp.Beers[bl.curBeer], nil
	}

	// otherwise we have to make a new request
	pageNum := bl.resp.CurrentPage + 1
	if pageNum > bl.resp.NumberOfPages {
		// no more pages
		return nil, nil
	}

	if err := bl.getPage(pageNum); err != nil {
		return nil, err
	}

	return &bl.resp.Beers[0], nil
}

// Get queries for a single Beer with the given Beer ID.
//
// TODO: add withBreweries, withSocialAccounts, withIngredients request parameters
func (bs *BeerService) Get(id string) (beer Beer, err error) {
	// GET: /beer/:beerId
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+id, nil)
	if err != nil {
		return
	}

	beerResp := struct {
		Message string
		Data    Beer
		Status  string
	}{}
	if err = bs.c.Do(req, &beerResp); err != nil {
		return
	}
	return beerResp.Data, nil
}

// BeerChangeRequest contains all the relevant options available to change
// an existing beer record in the BreweryDB.
type BeerChangeRequest struct {
	Name               string          `json:"name"`    // Required
	StyleID            int             `json:"styleId"` // Required
	Description        string          `json:"description"`
	ABV                string          `json:"abv"`
	IBU                string          `json:"ibu"`
	GlasswareID        int             `json:"glasswareId"`
	SrmID              int             `json:"srmID"`
	AvailableID        int             `json:"availableId"`
	IsOrganic          string          `json:"isOrganic"`
	BeerVariationID    string          `json:"beerVariationId"`
	Year               int             `json:"year"`
	FoodPairings       string          `json:"foodPairings"`
	ServingTemperature BeerTemperature `json:"servingTemperature"`
	OriginalGravity    string          `json:"originalGravity"`
	Brewery            string          `json:"brewery"` // Comma separated list of existing brewery IDs
	Label              string          `json:"label"`   // Base 64 encoded image
}

// Add adds a new Beer to the BreweryDB and returns its new ID on success.
func (bs *BeerService) Add(q *BeerChangeRequest) (id string, err error) {
	// POST: /beers
	var req *http.Request
	req, err = bs.c.NewRequest("POST", "/beers", q)
	if err != nil {
		return
	}

	addResponse := struct {
		Data struct {
			ID string
		}
	}{}
	if err = bs.c.Do(req, &addResponse); err != nil {
		return
	}

	return addResponse.Data.ID, nil
}

// Update changes an existing Beer in the BreweryDB.
func (bs *BeerService) Update(id string, q *BeerChangeRequest) error {
	// PUT: /beer/:beerId
	req, err := bs.c.NewRequest("PUT", "/beer/"+id, q)
	if err != nil {
		return err
	}

	// TODO: check status==success in JSON response body?
	return bs.c.Do(req, nil)
}

// Delete removes the Beer with the given ID from the BreweryDB.
func (bs *BeerService) Delete(id string) error {
	// DELETE: /beer/:beerId
	req, err := bs.c.NewRequest("DELETE", "/beer/"+id, nil)
	if err != nil {
		return err
	}

	// TODO: extract and return response message
	return bs.c.Do(req, nil)
}
