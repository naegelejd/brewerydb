package brewerydb

import "net/http"

// RandomBeerRequest contains options for querying for a random beer.
type RandomBeerRequest struct {
	ABV                string `json:"abv,omitempty"`
	IBU                string `json:"ibu,omitempty"`
	GlasswareID        int    `json:"glasswareId,omitempty"`
	SrmID              int    `json:"srmID,omitempty"`
	AvailableID        int    `json:"availableId,omitempty"`
	StyleID            int    `json:"styleId,omitempty"`
	IsOrganic          bool   `json:"isOrganic,omitempty"`
	Labels             bool   `json:"labels,omitempty"`
	Year               int    `json:"year,omitempty"`
	WithBreweries      bool   `json:"withBreweries,omitempty"`
	WithSocialAccounts bool   `json:"withSocialAccounts,omitempty"`
	WithIngredients    bool   `json:"withIngredients,omitempty"`
}

type randomBeerResponse struct {
	Status  string
	Message string
	Beer    Beer `json:"data"`
}

// Random returns a random beer that meets the requirements specified
// in the given RandomBeerRequest.
func (bs *BeerService) Random(q *RandomBeerRequest) (b Beer, err error) {
	// GET: /beer/random

	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/random", q)
	if err != nil {
		return
	}

	randomBeerResponse := struct {
		Status  string
		Data    Beer
		Message string
	}{}

	if err = bs.c.Do(req, &randomBeerResponse); err != nil {
		return
	}

	return randomBeerResponse.Data, nil
}
