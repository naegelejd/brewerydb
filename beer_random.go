package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RandomBeerRequest struct {
	ABV                string `json:"abv"`
	IBU                string `json:"ibu"`
	GlasswareId        string `json:"glasswareId"`
	SrmId              string `json:"srmID"`
	AvailableId        string `json:"availableId"`
	StyleId            string `json:"styleId"`
	IsOrganic          bool   `json:"isOrganic"`
	Labels             bool   `json:"labels"`
	Year               int    `json:"year"`
	WithBreweries      bool   `json:"withBreweries"`
	WithSocialAccounts bool   `json:"withSocialAccounts"`
	WithIngredients    bool   `json:"withIngredients"`
}

type randomBeerResponse struct {
	Status  string
	Message string
	Beer    Beer `json:"data"`
}

// GET: /beer/random
func (c *Client) RandomBeer(req *RandomBeerRequest) (b *Beer, err error) {
	vals := encode(req)

	u := c.URL("/beer/random", &vals)

	var resp *http.Response
	resp, err = c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	r := &randomBeerResponse{}
	if err = json.NewDecoder(resp.Body).Decode(r); err != nil {
		return
	}

	b = &r.Beer

	return
}
