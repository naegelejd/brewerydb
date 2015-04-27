package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type breweryListResponse struct {
	Status        string
	CurrentPage   int
	NumberOfPages int
	Breweries     []Brewery `json:"data"`
}

type BreweryList struct {
	c          *Client
	resp       *breweryListResponse
	req        *BreweryListRequest
	curBrewery int
}

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

// GET: /breweries
func (c *Client) NewBreweryList(req *BreweryListRequest) *BreweryList {
	return &BreweryList{c: c, req: req}
}

func (bl *BreweryList) getPage(pageNum int) error {
	v := encode(bl.req)
	v.Set("p", fmt.Sprintf("%d", pageNum))

	u := bl.c.url("/breweries", &v)

	resp, err := bl.c.Get(u)
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

// GET: /brewery/:breweryId
func (c *Client) Brewery(id string) (brewery *Brewery, err error) {
	u := c.url("/brewery/"+id, nil)
	var resp *http.Response
	resp, err = c.Get(u)
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

// POST: /breweries
func (c *Client) AddBrewery( /* params */ ) (id string, err error) {
	return
}

// PUT: /brewery/:breweryId
func (c *Client) UpdateBrewery( /* params */ ) error {
	return nil
}

// DELETE: /brewery/:breweryId
func (c *Client) DeleteBrewery( /* params */ ) error {
	return nil
}
