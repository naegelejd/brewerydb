package brewerydb

import (
	"fmt"
	"net/http"
)

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
	IsOrganic          YesNo     `url:"isOrganic,omitempty"`
	HasLabels          YesNo     `url:"hasLabels,omitempty"`
	Year               int       `url:"year,omitempty"`  // YYYY
	Since              string    `url:"since,omitempty"` // UNIX timestamp format. Max 30 days
	Status             string    `url:"status,omitempty"`
	Order              BeerOrder `url:"order,omitempty"`
	Sort               ListSort  `url:"sort,omitempty"`
	RandomCount        string    `url:"randomCount,omitempty"` // how many random beers to return. Max 10
	WithBreweries      YesNo     `url:"withBreweries,omitempty"`
	WithSocialAccounts YesNo     `url:"withSocialAccounts,omitempty"` // Premium
	WithIngredients    YesNo     `url:"withIngredients,omitempty"`    // Premium
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
	ID                        string          `url:"-"`
	Name                      string          `url:"name"` // Required
	Description               string          `url:"description,omitempty"`
	FoodPairings              string          `url:"foodPairings,omitempty"`
	OriginalGravity           string          `url:"originalGravity,omitempty"`
	ABV                       string          `url:"abv,omitempty"`
	IBU                       string          `url:"ibu,omitempty"`
	GlasswareID               int             `url:"glasswareId,omitempty"`
	Glass                     Glass           `url:"-"`
	StyleID                   int             `url:"styleId"` // Required
	Style                     Style           `url:"-"`
	IsOrganic                 YesNo           `url:"isOrganic,omitempty"`
	Labels                    Images          `url:"-"`
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
	SrmID      int    `url:"srmId,omitempty"`
	SRM        SRM    `url:"-"`
	Year       int    `url:"year,omitempty"`
	CreateDate string `url:"-"`
	UpdateDate string `url:"-"`
}

// List returns all Beers on the page specified in the given BeerListRequest.
// For non-premium members, one of Name, ABV, IBU, SrmID, AvailabilityID, StyleID must be set.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_index#1
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
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_index#2
func (bs *BeerService) Get(id string) (beer Beer, err error) {
	// GET: /beer/:beerId
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+id, nil)
	if err != nil {
		return
	}

	resp := struct {
		Message string
		Data    Beer
		Status  string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// Add adds a new Beer to the BreweryDB and returns its new ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_index#3
func (bs *BeerService) Add(b *Beer) (id string, err error) {
	// POST: /beers
	if b == nil {
		err = fmt.Errorf("nil Beer")
		return
	}
	var req *http.Request
	req, err = bs.c.NewRequest("POST", "/beers", b)
	if err != nil {
		return
	}

	resp := struct {
		Status string
		Data   struct {
			ID string
		}
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data.ID, err
}

// Update changes an existing Beer in the BreweryDB.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_index#4
func (bs *BeerService) Update(id string, b *Beer) error {
	// PUT: /beer/:beerId
	if b == nil {
		return fmt.Errorf("nil Beer")
	}
	req, err := bs.c.NewRequest("PUT", "/beer/"+id, b)
	if err != nil {
		return err
	}

	// TODO: check status==success in JSON response body?
	return bs.c.Do(req, nil)
}

// Delete removes the Beer with the given ID from the BreweryDB.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_index#5
func (bs *BeerService) Delete(id string) error {
	// DELETE: /beer/:beerId
	req, err := bs.c.NewRequest("DELETE", "/beer/"+id, nil)
	if err != nil {
		return err
	}

	// TODO: extract and return response message
	return bs.c.Do(req, nil)
}

// ListAdjuncts returns a slice of all Adjuncts in the Beer with the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_adjunct#1
func (bs *BeerService) ListAdjuncts(beerID string) (al []Adjunct, err error) {
	// GET: /beer/:beerId/adjuncts
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+beerID+"/adjuncts", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Adjunct
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// AddAdjunct adds the Adjunct with the given ID to the Beer with the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_adjunct#2
func (bs *BeerService) AddAdjunct(beerID string, adjunctID int) error {
	// POST: /beer/:beerId/adjuncts
	q := struct {
		ID int `url:"adjunctId"`
	}{adjunctID}
	req, err := bs.c.NewRequest("POST", "/beer/"+beerID+"/adjuncts", &q)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// DeleteAdjunct removes the Adjunct with the given ID from the Beer with the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_adjunct#3
func (bs *BeerService) DeleteAdjunct(beerID string, adjunctID int) error {
	// DELETE: /beer/:beerId/adjunct/:adjunctId
	req, err := bs.c.NewRequest("DELETE", fmt.Sprintf("/beer/%s/adjunct/%d", beerID, adjunctID), nil)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}

// ListBreweries queries for all Breweries associated with the Beer having the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_brewery#1
func (bs *BeerService) ListBreweries(id string) ([]Brewery, error) {
	// GET: /beer/:beerId/breweries
	req, err := bs.c.NewRequest("GET", "/beer/"+id+"/breweries", nil)
	if err != nil {
		return nil, err
	}

	resp := struct {
		Status  string
		Data    []Brewery
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// BeerBreweryRequest allows for specifying locations for a given Brewery, e.g.
// if only adding/removing a specific Brewery location from a Beer.
type BeerBreweryRequest struct {
	LocationID string `url:"locationId,omitempty"`
}

// AddBrewery adds the Brewery with the given Brewery ID to the Beer
// with the given Beer ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_brewery#2
func (bs *BeerService) AddBrewery(beerID, breweryID string, q *BeerBreweryRequest) error {
	// POST: /beer/:beerId/brewery/:breweryId
	params := struct {
		ID         string `url:"breweryId"`
		LocationID string `url:"locationId,omitempty"`
	}{ID: breweryID}

	if q != nil {
		params.LocationID = q.LocationID
	}

	req, err := bs.c.NewRequest("POST", "/beer/"+beerID+"/breweries", &params)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// DeleteBrewery removes the Brewery with the given Brewery ID from the Beer
// with the given Beer ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_brewery#3
func (bs *BeerService) DeleteBrewery(beerID, breweryID string, q *BeerBreweryRequest) error {
	// DELETE: /beer/:beerId/brewery/:breweryId
	req, err := bs.c.NewRequest("DELETE", "/beer/"+beerID+"/brewery/"+breweryID, q)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// ListEvents returns a slice of Events where the given Beer is/was present
// or has won awards.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_event#1
func (bs *BeerService) ListEvents(beerID string, onlyWinners bool) (el []Event, err error) {
	// GET: /beer/:beerId/events
	q := struct {
		OnlyWinners YesNo `url:"onlyWinners,omitempty"`
	}{YesNo(onlyWinners)}

	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+beerID+"/events", &q)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Event
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// ListFermentables returns a slice of all Fermentables in the Beer with the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_fermentable#1
func (bs *BeerService) ListFermentables(beerID string) (fl []Fermentable, err error) {
	// GET: /beer/:beerId/fermentables
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+beerID+"/fermentables", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Fermentable
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// AddFermentable adds the Fermentable with the given ID to the Beer with the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_fermentable#2
func (bs *BeerService) AddFermentable(beerID string, fermentableID int) error {
	// POST: /beer/:beerId/fermentables
	q := struct {
		ID int `url:"fermentableId"`
	}{fermentableID}
	req, err := bs.c.NewRequest("POST", "/beer/"+beerID+"/fermentables", &q)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}

// DeleteFermentable removes the Fermentable with the given ID from the Beer with the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_fermentable#3
func (bs *BeerService) DeleteFermentable(beerID string, fermentableID int) error {
	// DELETE: /beer/:beerId/fermentable/:fermentableId
	req, err := bs.c.NewRequest("DELETE", fmt.Sprintf("/beer/%s/fermentable/%d", beerID, fermentableID), nil)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}

// ListHops returns a slice of all Hops in the Beer with the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_hop#1
func (bs *BeerService) ListHops(beerID string) (al []Hop, err error) {
	// GET: /beer/:beerId/hops
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+beerID+"/hops", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Hop
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// AddHop adds the Hop with the given ID to the Beer with the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_hop#2
func (bs *BeerService) AddHop(beerID string, hopID int) error {
	// POST: /beer/:beerId/hops
	q := struct {
		ID int `url:"hopId"`
	}{hopID}
	req, err := bs.c.NewRequest("POST", "/beer/"+beerID+"/hops", &q)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// DeleteHop removes the Hop with the given ID from the Beer with the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_hop#3
func (bs *BeerService) DeleteHop(beerID string, hopID int) error {
	// DELETE: /beer/:beerId/hop/:hopId
	req, err := bs.c.NewRequest("DELETE", fmt.Sprintf("/beer/%s/hop/%d", beerID, hopID), nil)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}

// ListIngredients returns a slice of Ingredients found in the Beer with the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_ingredient#1
func (bs *BeerService) ListIngredients(beerID string) (el []Ingredient, err error) {
	// GET: /beer/:beerId/ingredients
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+beerID+"/ingredients", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Ingredient
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// RandomBeerRequest contains options for retrieving a random Beer.
type RandomBeerRequest struct {
	ABV                string `url:"abv,omitempty"`
	IBU                string `url:"ibu,omitempty"`
	GlasswareID        int    `url:"glasswareId,omitempty"`
	SrmID              int    `url:"srmID,omitempty"`
	AvailableID        int    `url:"availableId,omitempty"`
	StyleID            int    `url:"styleId,omitempty"`
	IsOrganic          YesNo  `url:"isOrganic,omitempty"`
	Labels             YesNo  `url:"labels,omitempty"`
	Year               int    `url:"year,omitempty"`
	WithBreweries      YesNo  `url:"withBreweries,omitempty"`
	WithSocialAccounts YesNo  `url:"withSocialAccounts,omitempty"`
	WithIngredients    YesNo  `url:"withIngredients,omitempty"`
}

// GetRandom returns a random Beer.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_random#1
func (bs *BeerService) GetRandom(q *RandomBeerRequest) (b Beer, err error) {
	// GET: /beer/random

	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/random", q)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    Beer
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, nil
}

// ListSocialAccounts returns a slice of all social media accounts associated with the given Beer.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_socialaccount#1
func (bs *BeerService) ListSocialAccounts(beerID string) (sl []SocialAccount, err error) {
	// GET: /beer/:beerId/socialaccounts
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+beerID+"/socialaccounts", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []SocialAccount
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// GetSocialAccount retrieves the SocialAccount with the given ID for the given Beer.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_socialaccount#2
func (bs *BeerService) GetSocialAccount(beerID string, socialAccountID int) (s SocialAccount, err error) {
	// GET: /beer/:beerId/socialaccount/:socialaccountId
	var req *http.Request
	req, err = bs.c.NewRequest("GET", fmt.Sprintf("/beer/%s/socialaccount/%d", beerID, socialAccountID), nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    SocialAccount
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// AddSocialAccount adds a new SocialAccount to the given Beer.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_socialaccount#3
func (bs *BeerService) AddSocialAccount(beerID string, s *SocialAccount) error {
	// POST: /beer/:beerId/socialaccounts
	if s == nil {
		return fmt.Errorf("nil SocialAccount")
	}
	req, err := bs.c.NewRequest("POST", "/beer/"+beerID+"/socialaccounts", s)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// UpdateSocialAccount updates a SocialAccount for the given Beer.
//
// See: TODO (API docs seem wrong)
func (bs *BeerService) UpdateSocialAccount(beerID string, s *SocialAccount) error {
	// PUT: /beer/:beerId/socialaccount/:socialaccountId
	if s == nil {
		return fmt.Errorf("nil SocialAccount")
	}
	req, err := bs.c.NewRequest("PUT", fmt.Sprintf("/beer/%s/socialaccount/%d", beerID, s.ID), s)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// DeleteSocialAccount removes a SocialAccount from the given Beer.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_socialaccount#4
func (bs *BeerService) DeleteSocialAccount(beerID string, socialAccountID int) error {
	// DELETE: /beer/:beerId/socialaccount/:socialaccountId
	req, err := bs.c.NewRequest("DELETE", fmt.Sprintf("/beer/%s/socialaccount/%d", beerID, socialAccountID), nil)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}

// AddUPC assigns a Universal Product Code to the Beer with the given ID.
// fluidsizeID is optional.
// NOTE: fluidsizeID is encoded as "fluidSizeId" with a capital S.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_upc#1
func (bs *BeerService) AddUPC(beerID string, code uint64, fluidsizeID *int) error {
	// POST: /beer/:beerId/upcs
	q := struct {
		Code        uint64 `url:"upcCode"`
		FluidsizeID int    `url:"fluidSizeId,omitempty"`
	}{Code: code}

	if fluidsizeID != nil {
		q.FluidsizeID = *fluidsizeID
	}

	req, err := bs.c.NewRequest("POST", "/beer/"+beerID+"/upcs", &q)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}

// ListVariations returns a slice of all Beers that are variations of the
// Beer with the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_variation#1
func (bs *BeerService) ListVariations(beerID string) (bl []Beer, err error) {
	// GET: /beer/:beerId/variations
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+beerID+"/variations", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Beer
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// ListYeasts returns a slice of all Yeasts in the Beer with the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_yeast#1
func (bs *BeerService) ListYeasts(beerID string) (al []Yeast, err error) {
	// GET: /beer/:beerId/yeasts
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+beerID+"/yeasts", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Yeast
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// AddYeast adds the Yeast with the given ID to the Beer with the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_yeast#2
func (bs *BeerService) AddYeast(beerID string, yeastID int) error {
	// POST: /beer/:beerId/yeasts
	q := struct {
		ID int `url:"yeastId"`
	}{yeastID}
	req, err := bs.c.NewRequest("POST", "/beer/"+beerID+"/yeasts", &q)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// DeleteYeast removes the Yeast with the given ID from the Beer with the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/beer_yeast#3
func (bs *BeerService) DeleteYeast(beerID string, yeastID int) error {
	// DELETE: /beer/:beerId/yeast/:yeastId
	req, err := bs.c.NewRequest("DELETE", fmt.Sprintf("/beer/%s/yeast/%d", beerID, yeastID), nil)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}
