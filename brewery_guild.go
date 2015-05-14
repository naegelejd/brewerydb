package brewerydb

import (
	"fmt"
	"net/http"
)

// ListGuilds returns a slice of all Guilds the Brewery with the given ID belongs to.
func (bs *BreweryService) ListGuilds(breweryID string) (al []Guild, err error) {
	// GET: /brewery/:breweryId/guilds
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/brewery/"+breweryID+"/guilds", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Guild
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// AddGuild adds the Guild with the given ID to the Brewery with the given ID.
// discount is optional (value of discount offered to guild members).
func (bs *BreweryService) AddGuild(breweryID string, guildID int, discount *string) error {
	// POST: /brewery/:breweryId/guilds
	q := struct {
		ID       int    `url:"guildId"`
		Discount string `url:"discount"`
	}{ID: guildID}
	if discount != nil {
		q.Discount = *discount
	}

	req, err := bs.c.NewRequest("POST", "/brewery/"+breweryID+"/guilds", &q)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// DeleteGuild removes the Guild with the given ID from the Brewery with the given ID.
func (bs *BreweryService) DeleteGuild(breweryID string, guildID int) error {
	// DELETE: /brewery/:breweryId/guild/:guildId
	req, err := bs.c.NewRequest("DELETE", fmt.Sprintf("/brewery/%s/guild/%d", breweryID, guildID), nil)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}
