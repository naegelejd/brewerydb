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
