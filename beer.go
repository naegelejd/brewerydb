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

// GET: /beer/:beerId/breweries
// POST: /beer/:beerId/breweries
// DELETE: /beer/:beerId/brewery/:breweryId

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

type BeerList struct {
	c       *Client
	req     *BeerListRequest
	resp    *beerListResponse
	curBeer int
}

type BeerListRequest struct {
}

type beerListResponse struct {
	Status        string
	CurrentPage   int
	NumberOfPages int
	Beers         []Beer `json:"data"`
}

type Beer struct {
	ID              string
	Name            string
	Description     string
	FoodPairings    string
	OriginalGravity string
	ABV             string
	IBU             string
	GlasswareID     float64
	Glass           struct {
		UpdateDate  string
		ID          float64
		Description string
		CreateDate  string
		Name        string
	}
	StyleID float64
	Style   struct {
		ID       float64
		Category struct {
			UpdateDate  string
			ID          float64
			Description float64
			CreateDate  string
			Name        string
		}
		Description string
		IbuMin      string
		IbuMax      string
		SrmMin      string
		SrmMax      string
		OgMin       string
		OgMax       string
		FgMin       string
		FgMax       string
		AbvMin      string
		AbvMax      string
		CreateDate  string
		UpdateDate  string
		Name        string
		CategoryID  float64
	}
	IsOrganic string
	Labels    struct {
		Medium string
		Large  string
		Icon   string
	}
	ServingTemperature        string
	ServingTemperatureDisplay string
	Status                    string
	StatusDisplay             string
	AvailableID               string
	Available                 struct {
		Description string
		Name        string
	}
	BeerVariationID string
	BeerVariation   struct {
		// TODO: instance of a Beer??
	}
	SrmID string
	Srm   struct {
		// TODO: empty array??
	}
	Year string
}

type beerResponse struct {
	Message string
	Beer    Beer `json:"data"`
	Status  string
}

// GET: /beers
func (c *Client) NewBeerList(req *BeerListRequest) *BeerList {
	return &BeerList{c: c, req: req}
}

func (bl *BeerList) getPage(pageNum int) error {
	v := url.Values{}
	v.Set("p", fmt.Sprintf("%d", pageNum))
	// TODO: encode bl.req (BeerListRequest)
	u := bl.c.URL("/beers", &v)

	resp, err := bl.c.Get(u)
	if err != nil {
		return err
	} else if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("beers not found")
	}
	defer resp.Body.Close()

	beerListResp := &beerListResponse{}
	if err := json.NewDecoder(resp.Body).Decode(beerListResp); err != nil {
		return err
	}

	if len(beerListResp.Beers) <= 0 {
		return fmt.Errorf("no beers found on page %d", pageNum)
	}

	bl.resp = beerListResp
	bl.curBeer = 0

	return nil
}

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

// GET: /beer/:beerId
func (c *Client) Beer(id string) (beer *Beer, err error) {
	u := c.URL("/beer/"+id, nil)
	var resp *http.Response
	resp, err = c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("beer not found")
		return
	}
	defer resp.Body.Close()

	beerResp := beerResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&beerResp); err != nil {
		return
	}
	beer = &beerResp.Beer
	return
}

// POST: /beers
func (c *Client) AddBeer( /* params */ ) (id string, err error) {
	return
}

// PUT: /beer/:beerId
func (c *Client) UpdateBeer( /* params */ ) error {
	return nil
}

// DELETE: /beer/:beerId
func (c *Client) DeleteBeer(id string) error {
	return nil
}
