package brewerydb

import "net/http"

// RandomBeerRequest contains options for retrieving a random Beer.
type RandomBeerRequest struct {
	ABV                string `url:"abv,omitempty"`
	IBU                string `url:"ibu,omitempty"`
	GlasswareID        int    `url:"glasswareId,omitempty"`
	SrmID              int    `url:"srmID,omitempty"`
	AvailableID        int    `url:"availableId,omitempty"`
	StyleID            int    `url:"styleId,omitempty"`
	IsOrganic          bool   `url:"isOrganic,omitempty"` // Y/N
	Labels             bool   `url:"labels,omitempty"`
	Year               int    `url:"year,omitempty"`
	WithBreweries      string `url:"withBreweries,omitempty"`      // Y/N
	WithSocialAccounts string `url:"withSocialAccounts,omitempty"` // Y/N
	WithIngredients    string `url:"withIngredients,omitempty"`    // Y/N
}

// Random returns a random Beer.
func (bs *BeerService) Random(q *RandomBeerRequest) (b Beer, err error) {
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
