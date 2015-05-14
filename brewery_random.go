package brewerydb

import "net/http"

// RandomBreweryRequest contains options for retrieving a random Brewery.
type RandomBreweryRequest struct {
	Established        string `url:"established"`        // YYYY
	IsOrganic          string `url:"isOrganic"`          // Y/N
	WithSocialAccounts string `url:"withSocialAccounts"` // Y/N
	WithGuilds         string `url:"withGuilds"`         // Y/N
	WithLocations      string `url:"withLocations"`      // Y/N
	WithAlternateNames string `url:"withAlternateNames"` // Y/N
}

// Random returns a random active Brewery.
func (bs *BreweryService) Random(q *RandomBreweryRequest) (b Brewery, err error) {
	// GET: /brewery/random
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/brewery/random", q)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    Brewery
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}
