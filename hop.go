package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type hopListResponse struct {
	Status        string
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Hops          []Hop `json:"data"`
}

type HopList struct {
	c      *Client
	resp   *hopListResponse
	curHop int
}

type Hop struct {
	ID               int
	Name             string
	Description      string
	CountryOfOrigin  string
	AlphaAcidMin     float64
	AlphaAcidMax     float64
	BetaAcidMin      float64
	BetaAcidMax      float64
	HumuleneMin      float64
	HumuleneMax      float64
	CaryophylleneMin float64
	CaryophylleneMax float64
	CohumuloneMin    float64
	CohumuloneMax    float64
	MyrceneMin       float64
	MyrceneMax       float64
	FarneseneMin     float64
	FarneseneMax     float64
	IsNoble          string
	ForBittering     string
	ForFlavor        string
	ForAroma         string
	Category         string
	CategoryDisplay  string
	CreateDate       string
	UpdateDate       string
	Country          struct {
		IsoCode     string
		Name        string
		DisplayName string
		IsoThree    string
		NumberCode  int
		CreateDate  string
	}
}

type hopResponse struct {
	Message string
	Hop     Hop `json:"data"`
	Status  string
}

// Hops returns a HopsList which can be consumed like so:
//
// hs, _ := client.NewHopList()
// for h, err := hs.First(); h != nil; h, err = hs.Next() {
//	if err != nil { ...; break }
//	...
// }
func (c *Client) NewHopList() *HopList {
	return &HopList{c: c}
}

func (hl *HopList) getPage(pageNum int) error {
	v := url.Values{}
	v.Set("p", fmt.Sprintf("%d", pageNum))
	u := hl.c.url("/hops", &v)

	resp, err := hl.c.Get(u)
	if err != nil {
		return err
	} else if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("hops not found")
	}
	defer resp.Body.Close()

	hopListResp := &hopListResponse{}
	if err := json.NewDecoder(resp.Body).Decode(hopListResp); err != nil {
		return err
	}

	if len(hopListResp.Hops) <= 0 {
		return fmt.Errorf("no hops found on page %d", pageNum)
	}

	hl.resp = hopListResp
	hl.curHop = 0

	return nil
}

// GET: /hops
func (hl *HopList) First() (*Hop, error) {
	// If we already have page 1 cached, just return the first Hop
	if hl.resp != nil && hl.resp.CurrentPage == 1 {
		hl.curHop = 0
		return &hl.resp.Hops[0], nil
	}

	if err := hl.getPage(1); err != nil {
		return nil, err
	}

	return &hl.resp.Hops[0], nil
}

func (hl *HopList) Next() (*Hop, error) {
	hl.curHop++
	// if we're still on the same page just return hop
	if hl.curHop < len(hl.resp.Hops) {
		return &hl.resp.Hops[hl.curHop], nil
	}

	// otherwise we have to make a new request
	pageNum := hl.resp.CurrentPage + 1
	if pageNum > hl.resp.NumberOfPages {
		// no more pages
		return nil, nil
	}

	if err := hl.getPage(pageNum); err != nil {
		return nil, err
	}

	return &hl.resp.Hops[0], nil
}

// GET: /hop/:hopId
func (c *Client) Hop(id int) (hop *Hop, err error) {
	u := c.url(fmt.Sprintf("/hop/%d", id), nil)
	var resp *http.Response
	resp, err = c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("hop not found")
		return
	}
	defer resp.Body.Close()

	hopResp := hopResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&hopResp); err != nil {
		return
	}
	hop = &hopResp.Hop
	return
}
