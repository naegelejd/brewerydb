package brewerydb

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

	u := c.url("/beer/random", &vals)

	r := &randomBeerResponse{}
	if err = c.getJSON(u, r); err != nil {
		return
	}

	b = &r.Beer

	return
}
