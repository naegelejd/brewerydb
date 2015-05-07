package brewerydb

import "net/http"

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
	ID             string
	Description    string
	Name           string
	CreateDate     string
	MailingListURL string
	UpdateDate     string
	Images         struct {
		Medium string
		Small  string
		Icon   string
	}
	Established   string
	IsOrganic     string
	Website       string
	Status        string
	StatusDisplay string
}

// BreweryListRequest contains all the required and optional fields
// used for querying for a list of Breweries.
type BreweryListRequest struct {
	Page        int          `json:"p"`
	Name        string       `json:"name,omitempty"`
	IDs         string       `json:"ids,omitempty"`
	Established string       `json:"established,omitempty"`
	IsOrganic   string       `json:"isOrganic,omitempty"`
	HasImages   string       `json:"hasImages,omitempty"`
	Since       string       `json:"since,omitempty"`
	Status      string       `json:"status,omitempty"`
	Order       BreweryOrder `json:"order,omitempty"` // TODO: enumerate
	Sort        string       `json:"sort,omitempty"`  // TODO: enumerate
	RandomCount string       `json:"randomCount,omitempty"`
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

// AddBrewery adds a new Brewery to the BreweryDB and returns its new ID on success.
// TODO: ensure a *Brewery can be used to add new brewery
func (bs *BreweryService) AddBrewery(b *Brewery) (id string, err error) {
	// POST: /breweries
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

// UpdateBrewery changes an existing Brewery in the BreweryDB.
// TODO: ensure a *Brewery can be used to add new brewery
func (bs *BreweryService) UpdateBrewery(breweryID string, b *Brewery) error {
	// PUT: /brewery/:breweryId
	req, err := bs.c.NewRequest("PUT", "/brewery/"+breweryID, b)
	if err != nil {
		return err
	}

	// TODO: extract and return response message?
	return bs.c.Do(req, nil)
}

// DeleteBrewery removes the Brewery with the given ID from the BreweryDB.
func (bs *BreweryService) DeleteBrewery(id string) error {
	// DELETE: /brewery/:breweryId
	req, err := bs.c.NewRequest("DELETE", "/brewery/"+id, nil)
	if err != nil {
		return err
	}

	// TODO: extract and return response message?
	return bs.c.Do(req, nil)
}
