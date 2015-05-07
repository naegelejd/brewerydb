package brewerydb

import "net/http"

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

// BeerList represents a "page" containing a slice of Beers.
type BeerList struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Beers         []Beer `json:"data"`
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
	Page               int       `json:"p"`
	IDs                string    `json:"ids,omitempty"` // IDs of the beers to return, comma separated. Max 10.
	Name               string    `json:"name,omitempty"`
	ABV                string    `json:"abv,omitempty"`
	IBU                string    `json:"ibu,omitempty"`
	GlasswareID        int       `json:"glasswareId,omitempty"`
	SrmID              int       `json:"srmId,omitempty"`
	AvailableID        int       `json:"availableId,omitempty"`
	StyleID            int       `json:"styleId,omitempty"`
	IsOrganic          string    `json:"isOrganic,omitempty"` // Y/N
	HasLabels          string    `json:"hasLabels,omitempty"` // Y/N
	Year               int       `json:"year,omitempty"`      // YYYY
	Since              string    `json:"since,omitempty"`     // UNIX timestamp format. Max 30 days
	Status             string    `json:"status,omitempty"`
	Order              BeerOrder `json:"order,omitempty"`
	Sort               ListSort  `json:"sort,omitempty"`
	RandomCount        string    `json:"randomCount,omitempty"`        // how many random beers to return. Max 10
	WithBreweries      string    `json:"withBreweries,omitempty"`      // Y/N
	WithSocialAccounts string    `json:"withSocialAccounts,omitempty"` // Premium. Y/N
	WithIngredients    string    `json:"withIngredients,omitempty"`    // Premium. Y/N
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

// List returns all Beers on the page specified in the given BeerListRequest.
func (bs *BeerService) List(q *BeerListRequest) (bl BeerList, err error) {
	// GET: /beers

	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beers", q)
	if err != nil {
		return
	}

	err = bs.c.Do(req, &bl)
	return
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
// TODO: remove this and just use type Beer
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
