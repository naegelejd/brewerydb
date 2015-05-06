package brewerydb

import "net/http"

// SearchService provides access to the BreweryDB Search API.
// Use Client.Search.
type SearchService struct {
	c *Client
}

type searchType string

const (
	searchBeer    searchType = "beer"
	searchBrewery            = "brewery"
	searchEvent              = "event"
	searchGuild              = "guild"
)

// SearchBeerResults is a list of Beers matching a Search.
// TODO: nearly identical to "BeerList" in beer.go
type SearchBeerResults struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Beers         []Beer `json:"data"`
}

// SearchBreweryResults is a list of Breweries matching a Search.
// TODO: nearly identical to "BreweryList" in beer.go
type SearchBreweryResults struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Beers         []Beer `json:"data"`
}

// SearchRequest contains options for narrowing a Search.
type SearchRequest struct {
	Page               int    `json:"p"`
	WithBreweries      string `json:"withBreweries,omitempty"`      // Y/N
	WithSocialAccounts string `json:"withSocialAccounts,omitempty"` // Y/N
	WithGuilds         string `json:"withGuilds,omitempty"`         // Y/N
	WithLocations      string `json:"withLocations,omitempty"`      // Y/N
	WithAlternateNames string `json:"withAlternateNames,omitempty"` // Y/N
	WithIngredients    string `json:"withIngredients,omitempty"`    // Y/N
}

type actualSearchRequest struct {
	Page               int        `json:"p"`
	Query              string     `json:"q"`
	Type               searchType `json:"type"`
	WithBreweries      string     `json:"withBreweries,omitempty"`
	WithSocialAccounts string     `json:"withSocialAccounts,omitempty"`
	WithGuilds         string     `json:"withGuilds,omitempty"`
	WithLocations      string     `json:"withLocations,omitempty"`
	WithAlternateNames string     `json:"withAlternateNames,omitempty"`
	WithIngredients    string     `json:"withIngredients,omitempty"`
}

func makeActualRequest(req *SearchRequest, query string, tp searchType) *actualSearchRequest {
	return &actualSearchRequest{
		Page:               req.Page,
		Query:              query,
		Type:               tp,
		WithBreweries:      req.WithBreweries,
		WithGuilds:         req.WithGuilds,
		WithLocations:      req.WithLocations,
		WithAlternateNames: req.WithAlternateNames,
		WithIngredients:    req.WithIngredients,
	}

}

// Beer searches for Beers matching the given query.
func (ss *SearchService) Beer(query string, q *SearchRequest) (sr SearchBeerResults, err error) {
	actualRequest := makeActualRequest(q, query, searchBeer)
	var req *http.Request
	req, err = ss.c.NewRequest("GET", "/search", actualRequest)
	if err != nil {
		return
	}

	err = ss.c.Do(req, &sr)
	return
}

// Brewery searches for Breweries matching the given query.
func (ss *SearchService) Brewery(query string, q *SearchRequest) (sr SearchBreweryResults, err error) {
	actualRequest := makeActualRequest(q, query, searchBrewery)
	var req *http.Request
	req, err = ss.c.NewRequest("GET", "/search", actualRequest)
	if err != nil {
		return
	}

	err = ss.c.Do(req, &sr)
	return
}

// Event searches for Events matching the given query.
func (ss *SearchService) Event(query string, q *SearchRequest) (el EventList, err error) {
	actualRequest := makeActualRequest(q, query, searchEvent)
	var req *http.Request
	req, err = ss.c.NewRequest("GET", "/search", actualRequest)
	if err != nil {
		return
	}

	err = ss.c.Do(req, &el)
	return
}

// Guild searches for Guilds matching the given query.
func (ss *SearchService) Guild(query string, q *SearchRequest) (gl GuildList, err error) {
	actualRequest := makeActualRequest(q, query, searchGuild)
	var req *http.Request
	req, err = ss.c.NewRequest("GET", "/search", actualRequest)
	if err != nil {
		return
	}

	err = ss.c.Do(req, &gl)
	return
}

// TODO: use this helper for all 4 search functions
func (ss *SearchService) search(asr *actualSearchRequest, data interface{}) error {
	req, err := ss.c.NewRequest("GET", "/search", asr)
	if err != nil {
		return err
	}
	return ss.c.Do(req, data)
}
