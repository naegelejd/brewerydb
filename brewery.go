package brewerydb

import (
	"fmt"
	"net/http"
)

// BreweryService provides access to the BreweryDB Brewery API. Use Client.Brewery.
type BreweryService struct {
	c *Client
}

// BreweryOrder represents the ordering of a list of Breweries.
type BreweryOrder string

// BreweryList ordering options.
const (
	BreweryOrderName           BreweryOrder = "name"
	BreweryOrderDescription                 = "description"
	BreweryOrderWebsite                     = "website"
	BreweryOrderEstablished                 = "established"
	BreweryOrderMailingListURL              = "mailingListUrl"
	BreweryOrderIsOrganic                   = "isOrganic"
	BreweryOrderStatus                      = "status"
	BreweryOrderCreateDate                  = "createDate"
	BreweryOrderUpdateDate                  = "updateDate"
	BreweryOrderRandom                      = "random"
)

// BreweryList represents a "page" containing one slice of Breweries.
type BreweryList struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Breweries     []Brewery `json:"data"`
}

// Brewery contains all relevant information for a single Brewery.
type Brewery struct {
	ID             string `url:"-"`
	Name           string `url:"name"`
	Description    string `url:"description,omitempty"`
	MailingListURL string `url:"mailingListUrl,omitempty"`
	Images         struct {
		Medium string `url:"-"`
		Small  string `url:"-"`
		Icon   string `url:"-"`
	} `url:"-"`
	Image         string `url:"image,omitempty"` // only used for adding/update Breweries
	Established   string `url:"established,omitempty"`
	IsOrganic     string `url:"isOrganic,omitempty"`
	Website       string `url:"website,omitempty"`
	Status        string `url:"-"`
	StatusDisplay string `url:"-"`
	CreateDate    string `url:"-"`
	UpdateDate    string `url:"-"`
}

// BreweryListRequest contains all the required and optional fields
// used for querying for a list of Breweries.
type BreweryListRequest struct {
	Page        int          `url:"p"`
	Name        string       `url:"name,omitempty"`
	IDs         string       `url:"ids,omitempty"`
	Established string       `url:"established,omitempty"`
	IsOrganic   string       `url:"isOrganic,omitempty"`
	HasImages   string       `url:"hasImages,omitempty"`
	Since       string       `url:"since,omitempty"`
	Status      string       `url:"status,omitempty"`
	Order       BreweryOrder `url:"order,omitempty"` // TODO: enumerate
	Sort        string       `url:"sort,omitempty"`  // TODO: enumerate
	RandomCount string       `url:"randomCount,omitempty"`
	// TODO: premium account parameters
}

// List returns all Breweries on the page specified in the given BreweryListRequest.
func (bs *BreweryService) List(q *BreweryListRequest) (bl BreweryList, err error) {
	// GET: /breweries
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/breweries", q)
	if err != nil {
		return
	}

	err = bs.c.Do(req, &bl)
	return
}

// Get queries for a single Brewery with the given Brewery ID.
func (bs *BreweryService) Get(id string) (brewery Brewery, err error) {
	// GET: /brewery/:breweryId
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/brewery/"+id, nil)
	if err != nil {
		return
	}

	breweryResp := struct {
		Message string
		Data    Brewery
		Status  string
	}{}
	if err = bs.c.Do(req, &breweryResp); err != nil {
		return
	}
	return breweryResp.Data, nil
}

// Add adds a new Brewery to the BreweryDB and returns its new ID on success.
func (bs *BreweryService) Add(b *Brewery) (id string, err error) {
	// POST: /breweries
	if b == nil {
		err = fmt.Errorf("nil Brewery")
		return
	}
	var req *http.Request
	req, err = bs.c.NewRequest("POST", "/breweries", b)
	if err != nil {
		return
	}

	addResp := struct {
		Status  string
		ID      string
		Message string
	}{}
	err = bs.c.Do(req, &addResp)
	return addResp.ID, err
}

// Update changes an existing Brewery in the BreweryDB.
func (bs *BreweryService) Update(breweryID string, b *Brewery) error {
	// PUT: /brewery/:breweryId
	if b == nil {
		return fmt.Errorf("nil Brewery")
	}
	req, err := bs.c.NewRequest("PUT", "/brewery/"+breweryID, b)
	if err != nil {
		return err
	}

	// TODO: extract and return response message?
	return bs.c.Do(req, nil)
}

// Delete removes the Brewery with the given ID from the BreweryDB.
func (bs *BreweryService) Delete(id string) error {
	// DELETE: /brewery/:breweryId
	req, err := bs.c.NewRequest("DELETE", "/brewery/"+id, nil)
	if err != nil {
		return err
	}

	// TODO: extract and return response message?
	return bs.c.Do(req, nil)
}

// AlternateName represents an alternate name for a Brewery.
// TODO: the actual response object contains the entire Brewery object as well.
//	see: http://www.brewerydb.com/developers/docs-endpoint/brewery_alternatename
type AlternateName struct {
	ID         int
	Name       string
	BreweryID  string
	CreateDate string
	UpdateDate string
}

// ListAlternateNames returns a slice of all the AlternateNames for the Brewery with the given ID.
func (bs *BreweryService) ListAlternateNames(breweryID string) (al []AlternateName, err error) {
	// GET: /brewery/:breweryId/alternatenames
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/brewery/"+breweryID+"/alternatenames", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []AlternateName
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// AddAlternateName adds an alternate name to the Brewery with the given ID.
func (bs *BreweryService) AddAlternateName(breweryID, name string) error {
	// POST: /brewery/:breweryId/alternatenames
	q := struct {
		Name string `url:"name"`
	}{name}
	req, err := bs.c.NewRequest("POST", "/brewery/"+breweryID+"/alternatenames", &q)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}

// DeleteAlternateName removes the AlternateName with the given ID from the Brewery with the given ID.
func (bs *BreweryService) DeleteAlternateName(breweryID string, alternateNameID int) error {
	// DELETE: /brewery/:breweryId/alternatename/:alternatenameId
	req, err := bs.c.NewRequest("DELETE", fmt.Sprintf("/brewery/%s/alternatename/%d", breweryID, alternateNameID), nil)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}

// BreweryBeersRequest contains options for querying for all Beers from a Brewery.
type BreweryBeersRequest struct {
	WithBreweries      string `url:"withBreweries"`      // Y/N
	WithSocialAccounts string `url:"withSocialAccounts"` // Y/N
	WithIngredients    string `url:"withIngredients"`    // Y/N
}

// ListBeers returns a slice of all Beers offered by the Brewery with the given ID.
func (bs *BreweryService) ListBeers(breweryID string, q *BreweryBeersRequest) (bl []Beer, err error) {
	// GET: /brewery/:breweryId/beers
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/brewery/"+breweryID+"/beers", q)
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

// ListEvents returns a slice of Events where the given Brewery is/was present
// or has won awards.
func (bs *BreweryService) ListEvents(breweryID string, onlyWinners bool) (el []Event, err error) {
	// GET: /brewery/:breweryId/events
	var q struct {
		OnlyWinners string `url:"onlyWinners,omitempty"`
	}
	if onlyWinners {
		q.OnlyWinners = "Y"
	}

	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/brewery/"+breweryID+"/events", &q)
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

// ListGuilds returns a slice of all Guilds the Brewery with the given ID belongs to.
func (bs *BreweryService) ListGuilds(breweryID string) (al []Guild, err error) {
	// GET: /brewery/:breweryId/guilds
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/brewery/"+breweryID+"/guilds", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Guild
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// AddGuild adds the Guild with the given ID to the Brewery with the given ID.
// discount is optional (value of discount offered to guild members).
func (bs *BreweryService) AddGuild(breweryID string, guildID string, discount *string) error {
	// POST: /brewery/:breweryId/guilds
	q := struct {
		ID       string `url:"guildId"`
		Discount string `url:"discount"`
	}{ID: guildID}
	if discount != nil {
		q.Discount = *discount
	}

	req, err := bs.c.NewRequest("POST", "/brewery/"+breweryID+"/guilds", &q)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// DeleteGuild removes the Guild with the given ID from the Brewery with the given ID.
func (bs *BreweryService) DeleteGuild(breweryID string, guildID string) error {
	// DELETE: /brewery/:breweryId/guild/:guildId
	req, err := bs.c.NewRequest("DELETE", "/brewery/"+breweryID+"/guild/"+guildID, nil)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}

// ListLocations returns a slice of all locations for the Brewery with the given ID.
func (bs *BreweryService) ListLocations(breweryID string) (ll []Location, err error) {
	// GET: /brewery/:breweryId/locations
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/brewery/"+breweryID+"/locations", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Location
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// AddLocation adds a new location for the Brewery with the given ID.
func (bs *BreweryService) AddLocation(breweryID string, loc *Location) error {
	// POST: /brewery/:breweryId/locations
	if loc == nil {
		return fmt.Errorf("nil Location")
	}
	req, err := bs.c.NewRequest("POST", "/brewery/"+breweryID+"/locations", loc)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// RandomBreweryRequest contains options for retrieving a random Brewery.
type RandomBreweryRequest struct {
	Established        string `url:"established"`        // YYYY
	IsOrganic          string `url:"isOrganic"`          // Y/N
	WithSocialAccounts string `url:"withSocialAccounts"` // Y/N
	WithGuilds         string `url:"withGuilds"`         // Y/N
	WithLocations      string `url:"withLocations"`      // Y/N
	WithAlternateNames string `url:"withAlternateNames"` // Y/N
}

// GetRandom returns a random active Brewery.
func (bs *BreweryService) GetRandom(q *RandomBreweryRequest) (b Brewery, err error) {
	// GET: /brewery/random
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/brewery/random", q)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    Brewery
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// ListSocialAccounts returns a slice of all social media accounts associated with the given Brewery.
func (bs *BreweryService) ListSocialAccounts(breweryID string) (sl []SocialAccount, err error) {
	// GET: /brewery/:breweryId/socialaccounts
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/brewery/"+breweryID+"/socialaccounts", nil)
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

// GetSocialAccount retrieves the SocialAccount with the given ID for the given Brewery.
func (bs *BreweryService) GetSocialAccount(breweryID string, socialAccountID int) (s SocialAccount, err error) {
	// GET: /brewery/:breweryId/socialaccount/:socialaccountId
	var req *http.Request
	req, err = bs.c.NewRequest("GET", fmt.Sprintf("/brewery/%s/socialaccount/%d", breweryID, socialAccountID), nil)
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

// AddSocialAccount adds a new SocialAccount to the given Brewery.
func (bs *BreweryService) AddSocialAccount(breweryID string, s *SocialAccount) error {
	// POST: /brewery/:breweryId/socialaccounts
	if s == nil {
		return fmt.Errorf("nil SocialAccount")
	}
	req, err := bs.c.NewRequest("POST", "/brewery/"+breweryID+"/socialaccounts", s)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// UpdateSocialAccount updates a SocialAccount for the given Brewery.
func (bs *BreweryService) UpdateSocialAccount(breweryID string, s *SocialAccount) error {
	// PUT: /brewery/:breweryId/socialaccount/:socialaccountId
	if s == nil {
		return fmt.Errorf("nil SocialAccount")
	}
	req, err := bs.c.NewRequest("PUT", fmt.Sprintf("/brewery/%s/socialaccount/%d", breweryID, s.ID), s)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// DeleteSocialAccount removes a SocialAccount from the given Brewery.
func (bs *BreweryService) DeleteSocialAccount(breweryID string, socialAccountID int) error {
	// DELETE: /brewery/:breweryId/socialaccount/:socialaccountId
	req, err := bs.c.NewRequest("DELETE", fmt.Sprintf("/brewery/%s/socialaccount/%d", breweryID, socialAccountID), nil)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}
