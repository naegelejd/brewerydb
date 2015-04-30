package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// BreweryService provides access to the BreweryDB Brewery API. Use Client.Brewery.
type BreweryService struct {
	c *Client
}

type breweryListResponse struct {
	Status        string
	CurrentPage   int
	NumberOfPages int
	Breweries     []Brewery `json:"data"`
}

// BreweryList represents a lazy list of breweries. Create a new one with
// NewBreweryList. Iterate over a BreweryList using First() and Next().
type BreweryList struct {
	service    *BreweryService
	resp       *breweryListResponse
	req        *BreweryListRequest
	curBrewery int
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
	Name        string `json:"name"`
	IDs         string `json:"ids"`
	Established string `json:"established"`
	IsOrganic   string `json:"isOrganic"`
	HasImages   string `json:"hasImages"`
	Since       string `json:"since"`
	Status      string `json:"status"`
	Order       string `json:"order"` // TODO: enumerate
	Sort        string `json:"sort"`  // TODO: enumerate
	RandomCount string `json:"randomCount"`
	// TODO: premium account parameters
}

type breweryResponse struct {
	Message string
	Brewery Brewery `json:"data"`
	Status  string
}

// NewBreweryList returns a new BreweryList that will use the given
// BreweryListRequest to query for a list of Breweries.
func (s *BreweryService) NewBreweryList(req *BreweryListRequest) *BreweryList {
	// GET: /breweries
	return &BreweryList{service: s, req: req}
}

// getPage obtains the "next" page from the BreweryDB API
func (bl *BreweryList) getPage(pageNum int) error {
	v := encode(bl.req)
	v.Set("p", fmt.Sprintf("%d", pageNum))

	u := bl.service.c.url("/breweries", &v)

	resp, err := bl.service.c.Get(u)
	if err != nil {
		return err
	} else if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("breweries not found")
	}
	defer resp.Body.Close()

	breweryListResp := &breweryListResponse{}
	if err := json.NewDecoder(resp.Body).Decode(breweryListResp); err != nil {
		return err
	}

	if len(breweryListResp.Breweries) <= 0 {
		return fmt.Errorf("no breweries found on page %d", pageNum)
	}

	bl.resp = breweryListResp
	bl.curBrewery = 0

	return nil
}

// First returns the first Brewery in the BreweryList.
func (bl *BreweryList) First() (*Brewery, error) {
	// If we already have page 1 cached, just return the first Brewery
	if bl.resp != nil && bl.resp.CurrentPage == 1 {
		bl.curBrewery = 0
		return &bl.resp.Breweries[0], nil
	}

	if err := bl.getPage(1); err != nil {
		return nil, err
	}

	return &bl.resp.Breweries[0], nil
}

// Next returns the next Brewery in the BreweryList on each successive call,
// or nil if there are no more Breweries.
func (bl *BreweryList) Next() (*Brewery, error) {
	bl.curBrewery++
	// if we're still on the same page just return brewery
	if bl.curBrewery < len(bl.resp.Breweries) {
		return &bl.resp.Breweries[bl.curBrewery], nil
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

	return &bl.resp.Breweries[0], nil
}

// Get queries for a single Brewery with the given Brewery ID.
func (s *BreweryService) Get(id string) (brewery *Brewery, err error) {
	// GET: /brewery/:breweryId
	u := s.c.url("/brewery/"+id, nil)
	var resp *http.Response
	resp, err = s.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("brewery not found")
		return
	}
	defer resp.Body.Close()

	breweryResp := breweryResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&breweryResp); err != nil {
		return
	}
	brewery = &breweryResp.Brewery
	return
}

// AddBrewery adds a new Brewery to the BreweryDB and returns its new ID on success.
func (s *BreweryService) AddBrewery( /* params */ ) (id string, err error) {
	// POST: /breweries
	return
}

// UpdateBrewery changes an existing Brewery in the BreweryDB.
func (s *BreweryService) UpdateBrewery( /* params */ ) error {
	// PUT: /brewery/:breweryId
	return nil
}

// DeleteBrewery removes the Brewery with the given ID from the BreweryDB.
func (s *BreweryService) DeleteBrewery(id string) error {
	// DELETE: /brewery/:breweryId
	u := s.c.url("/brewery/"+id, nil)
	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	resp, err := s.c.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("brewery not found")
	}

	defer resp.Body.Close()

	// TODO: extract and return response message

	return nil
}
