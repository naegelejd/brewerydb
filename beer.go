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
	Page               int       `url:"p"`
	IDs                string    `url:"ids,omitempty"` // IDs of the beers to return, comma separated. Max 10.
	Name               string    `url:"name,omitempty"`
	ABV                string    `url:"abv,omitempty"`
	IBU                string    `url:"ibu,omitempty"`
	GlasswareID        int       `url:"glasswareId,omitempty"`
	SrmID              int       `url:"srmId,omitempty"`
	AvailableID        int       `url:"availableId,omitempty"`
	StyleID            int       `url:"styleId,omitempty"`
	IsOrganic          string    `url:"isOrganic,omitempty"` // Y/N
	HasLabels          string    `url:"hasLabels,omitempty"` // Y/N
	Year               int       `url:"year,omitempty"`      // YYYY
	Since              string    `url:"since,omitempty"`     // UNIX timestamp format. Max 30 days
	Status             string    `url:"status,omitempty"`
	Order              BeerOrder `url:"order,omitempty"`
	Sort               ListSort  `url:"sort,omitempty"`
	RandomCount        string    `url:"randomCount,omitempty"`        // how many random beers to return. Max 10
	WithBreweries      string    `url:"withBreweries,omitempty"`      // Y/N
	WithSocialAccounts string    `url:"withSocialAccounts,omitempty"` // Premium. Y/N
	WithIngredients    string    `url:"withIngredients,omitempty"`    // Premium. Y/N
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
	ID              string `url:"-"`
	Name            string `url:"name"` // Required
	Description     string `url:"description,omitempty"`
	FoodPairings    string `url:"foodPairings,omitempty"`
	OriginalGravity string `url:"originalGravity,omitempty"`
	ABV             string `url:"abv,omitempty"`
	IBU             string `url:"ibu,omitempty"`
	GlasswareID     int    `url:"glasswareId,omitempty"`
	Glass           Glass  `url:"-"`
	StyleID         int    `url:"styleId"` // Required
	Style           Style  `url:"-"`
	IsOrganic       string `url:"isOrganic,omitempty"`
	Labels          struct {
		Medium string `url:"-"`
		Large  string `url:"-"`
		Icon   string `url:"-"`
	} `url:"-"`
	Label                     string          `url:"label,omitempty"`   // base64. Only used for adding/updating Beers.
	Brewery                   []string        `url:"brewery,omitempty"` // breweryID list. Only used for adding/updating Beers.
	ServingTemperature        BeerTemperature `url:"servingTemperature,omitempty"`
	ServingTemperatureDisplay string          `url:"-"`
	Status                    string          `url:"-"`
	StatusDisplay             string          `url:"-"`
	AvailableID               int             `url:"availableId,omitempty"`
	Available                 Availability    `url:"-"`
	BeerVariationID           string          `url:"beerVariationId,omitempty"`
	BeerVariation             struct {
		// TODO: instance of a Beer??
	} `url:"-"`
	SrmID      int    `url:"srmID,omitempty"`
	SRM        SRM    `url:"-"`
	Year       int    `url:"year,omitempty"`
	CreateDate string `url:"-"`
	UpdateDate string `url:"-"`
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

// Add adds a new Beer to the BreweryDB and returns its new ID on success.
func (bs *BeerService) Add(b *Beer) (id string, err error) {
	// POST: /beers
	var req *http.Request
	req, err = bs.c.NewRequest("POST", "/beers", b)
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
func (bs *BeerService) Update(id string, b *Beer) error {
	// PUT: /beer/:beerId
	req, err := bs.c.NewRequest("PUT", "/beer/"+id, b)
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
