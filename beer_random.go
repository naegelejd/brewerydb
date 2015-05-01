package brewerydb

// RandomBeerRequest contains options for querying for a random beer.
type RandomBeerRequest struct {
	ABV                string `json:"abv"`
	IBU                string `json:"ibu"`
	GlasswareID        string `json:"glasswareId"`
	SrmID              string `json:"srmID"`
	AvailableID        string `json:"availableId"`
	StyleID            string `json:"styleId"`
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

// Random returns a random beer that meets the requirements specified
// in the given RandomBeerRequest.
func (s *BeerService) Random(req *RandomBeerRequest) (b *Beer, err error) {
	// GET: /beer/random
	vals := encode(req)

	u := s.c.url("/beer/random", &vals)

	r := &randomBeerResponse{}
	if err = s.c.getJSON(u, r); err != nil {
		return
	}

	b = &r.Beer

	return
}
