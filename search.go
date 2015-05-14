package brewerydb

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

func makeActualSearchRequest(req *SearchRequest, query string, tp searchType) *actualSearchRequest {
	if req == nil {
		return &actualSearchRequest{Query: query, Type: tp}
	}
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
	err = ss.search(query, q, searchBeer, &bl)
	return
}

// Brewery searches for Breweries matching the given query.
func (ss *SearchService) Brewery(query string, q *SearchRequest) (bl BreweryList, err error) {
	err = ss.search(query, q, searchBrewery, &bl)
	return
}

// Event searches for Events matching the given query.
func (ss *SearchService) Event(query string, q *SearchRequest) (el EventList, err error) {
	err = ss.search(query, q, searchEvent, &el)
	return
}

// Guild searches for Guilds matching the given query.
func (ss *SearchService) Guild(query string, q *SearchRequest) (gl GuildList, err error) {
	err = ss.search(query, q, searchGuild, &gl)
	return
}

func (ss *SearchService) search(query string, q *SearchRequest, t searchType, data interface{}) error {
	asr := makeActualSearchRequest(q, query, t)
	req, err := ss.c.NewRequest("GET", "/search", asr)
	if err != nil {
		return err
	}
	return ss.c.Do(req, data)
}
