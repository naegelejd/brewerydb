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

// SearchRequest contains options for narrowing a Search.
type SearchRequest struct {
	Page               int
	WithBreweries      string
	WithSocialAccounts string
	WithGuilds         string
	WithLocations      string
	WithAlternateNames string
	WithIngredients    string
}

type actualSearchRequest struct {
	Page               int        `url:"p,omitempty"`
	Query              string     `url:"q"` // required
	Type               searchType `url:"type"`
	WithBreweries      string     `url:"withBreweries,omitempty"`
	WithSocialAccounts string     `url:"withSocialAccounts,omitempty"`
	WithGuilds         string     `url:"withGuilds,omitempty"`
	WithLocations      string     `url:"withLocations,omitempty"`
	WithAlternateNames string     `url:"withAlternateNames,omitempty"`
	WithIngredients    string     `url:"withIngredients,omitempty"`
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
func (ss *SearchService) Beer(query string, q *SearchRequest) (bl BeerList, err error) {
	actualRequest := makeActualRequest(q, query, searchBeer)
	var req *http.Request
	req, err = ss.c.NewRequest("GET", "/search", actualRequest)
	if err != nil {
		return
	}

	err = ss.c.Do(req, &bl)
	return
}

// Brewery searches for Breweries matching the given query.
func (ss *SearchService) Brewery(query string, q *SearchRequest) (bl BreweryList, err error) {
	actualRequest := makeActualRequest(q, query, searchBrewery)
	var req *http.Request
	req, err = ss.c.NewRequest("GET", "/search", actualRequest)
	if err != nil {
		return
	}

	err = ss.c.Do(req, &bl)
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
