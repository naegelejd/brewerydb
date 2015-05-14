package brewerydb

import "net/http"

// ListBreweries returns a slice of all Breweries that are members of the given Guild.
func (gs *GuildService) ListBreweries(guildID string) (bl []Brewery, err error) {
	// GET: /guild/:guildId/breweries
	var req *http.Request
	req, err = gs.c.NewRequest("GET", "/guild/"+guildID+"/breweries", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Brewery
		Message string
	}{}
	err = gs.c.Do(req, &resp)
	return resp.Data, err
}
