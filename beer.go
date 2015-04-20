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

type ListOrder string
type ListSort string

const (
	NameOrder        ListOrder = "name"
	DescriptionOrder           = "description"
	AbvOrder                   = "abv"
	IbuOrder                   = "ibu"
	GlasswareIDOrder           = "glasswareId"
	SrmIDOrder                 = "smrID"
	AvailableIDOrder           = "availableId"
	StyleIDOrder               = "styleId"
	IsOrganicOrder             = "isOrganic"
	StatusOrder                = "status"
	CreateDateOrder            = "createDate"
	UpdateDateOrder            = "updateDate"
	RandomOrder                = "random"
)

const (
	AscendingSort  ListSort = "ASC"
	DescendingSort          = "DESC"
)

type BeerListRequest struct {
	IDs                string `json:"ids"` // IDs of the beers to return, comma separated. Max 10.
	Name               string `json:"name"`
	ABV                string `json:"abv"`
	IBU                string `json:"ibu"`
	GlasswareId        string `json:"glasswareId"`
	SrmId              string `json:"srmId"`
	AvailableId        string `json:"availableId"`
	StyleId            string `json:"styleId"`
	IsOrganic          string `json:"isOrganic"` // Y/N
	HasLabels          string `json:"hasLabels"` // Y/N
	Year               string `json:"year"`      // YYYY
	Since              string `json:"since"`     // UNIX timestamp format. Max 30 days
	Status             string `json:"status"`
	Order              string `json:"order"`
	Sort               string `json:"sort"`
	RandomCount        string `json:"randomCount"`        // how many random beers to return. Max 10
	WithBreweries      string `json:"withBreweries"`      // Y/N
	WithSocialAccounts string `json:"withSocialAccounts"` // Premium. Y/N
	WithIngredients    string `json:"withIngredients"`    // Premium. Y/N
}

type beerListResponse struct {
	Status        string
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
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
		ID          float64
		Name        string
		Description string
		CreateDate  string
		UpdateDate  string
	}
	StyleID float64
	Style   struct {
		ID         float64
		CategoryID float64
		Category   struct {
			ID          float64
			Name        string
			Description string
			CreateDate  string
			UpdateDate  string
		}
		Name        string
		ShortName   string
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
	AvailableID               float64
	Available                 struct {
		ID          float64
		Name        string
		Description string
	}
	BeerVariationID string
	BeerVariation   struct {
		// TODO: instance of a Beer??
	}
	SrmID float64
	Srm   struct {
		ID   float64
		Name string
		Hex  string
		// TODO: anything else?
	}
	Year string
}

// GET: /beers
func (c *Client) NewBeerList(req *BeerListRequest) *BeerList {
	return &BeerList{c: c, req: req}
}

func (bl *BeerList) getPage(pageNum int) error {
	var v url.Values
	if bl.req != nil {
		v = encode(bl.req)
	} else {
		v = url.Values{}
	}
	v.Set("p", fmt.Sprintf("%d", pageNum))

	u := bl.c.url("/beers", &v)

	resp, err := bl.c.Get(u)
	if err != nil {
		return err
	} else if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("beers not found")
	}
	defer resp.Body.Close()

	beerListResp := &beerListResponse{}
	if err := json.NewDecoder(resp.Body).Decode(beerListResp); err != nil {
		// if e, ok := err.(*json.UnmarshalTypeError); ok == true {
		// 	fmt.Printf("(JSON error) Value: %s, Type: %s", e.Value, e.Type.Kind())
		// }
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
// TODO: add withBreweries, withSocialAccounts, withIngredients request parameters
func (c *Client) Beer(id string) (beer *Beer, err error) {
	u := c.url("/beer/"+id, nil)
	var resp *http.Response
	resp, err = c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("beer not found")
		return
	}
	defer resp.Body.Close()

	beerResp := struct {
		Message string
		Beer    Beer `json:"data"`
		Status  string
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&beerResp); err != nil {
		return
	}
	beer = &beerResp.Beer
	return
}

type BeerChangeRequest struct {
	Name               string `json:"name"`    // Required
	StyleId            int    `json:"styleId"` // Required
	Description        string `json:"description"`
	ABV                string `json:"abv"`
	IBU                string `json:"ibu"`
	GlasswareId        int    `json:"glasswareId"`
	SrmId              string `json:"srmID"`
	AvailableId        string `json:"availableId"`
	IsOrganic          string `json:"isOrganic"`
	BeerVariationId    string `json:"beerVariationId"`
	Year               string `json:"year"`
	FoodPairings       string `json:"foodPairings"`
	ServingTemperature string `json:"servingTemperature"`
	OriginalGravity    string `json:"originalGravity"`
	Brewery            string `json:"brewery"` // Comma separated list of existing brewery IDs
	Label              string `json:"label"`   // Base 64 encoded image
}

// POST: /beers
func (c *Client) AddBeer(req *BeerChangeRequest) (id string, err error) {
	u := c.url("/beers", nil)

	var resp *http.Response
	resp, err = c.PostForm(u, encode(req))
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to create beer")
		return
	}
	defer resp.Body.Close()

	out := struct{ Data struct{ ID string } }{}
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return
	}

	id = out.Data.ID
	return
}

// PUT: /beer/:beerId
func (c *Client) UpdateBeer(id string, req *BeerChangeRequest) error {
	u := c.url("/beer/"+id, nil)

	resp, err := c.PostForm(u, encode(req))
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to update beer")
	}
	defer resp.Body.Close()

	// TODO: check "status"=="success" in JSON response body?

	return nil
}

// DELETE: /beer/:beerId
func (c *Client) DeleteBeer(id string) error {
	u := c.url("/beer/"+id, nil)
	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("beer not found")
	}

	defer resp.Body.Close()

	// TODO: Move to unit test and mock
	// m := make(map[string]string)
	// if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
	// 	return err
	// }
	// if m["status"] != "success" {
	// 	return fmt.Errorf("delete unsuccessful")
	// }

	return nil
}
