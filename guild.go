package brewerydb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GuildService provides access to the BreweryDB Guild API.
// Use Client.Guild.
type GuildService struct {
	c *Client
}

// Guild represents a Beer or Brewing organization.
type Guild struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"` // Required
	Description string `json:"description,omitempty"`
	Website     string `json:"website,omitempty"`
	Image       string `json:"image,omitempty"` // Base64
	Images      struct {
		Icon   string `json:"icon,omitempty"`
		Medium string `json:"medium,omitempty"`
		Large  string `json:"large,omitempty"`
	} `json:"Images,omitempty"`
	Established int `json:"established,omitempty"`
}

// GuildOrder specifies ordering of a GuildList.
type GuildOrder string

// GuildList ordering options.
const (
	GuildOrderName        GuildOrder = "name"
	GuildOrderDescription            = "description"
	GuildOrderEstablished            = "established"
	GuildOrderStatus                 = "status"
	GuildOrderCreateDate             = "createDate"
	GuildOrderUpdateDate             = "updateDate"
)

// GuildRequest contains options for specifying kinds of Guilds desired.
type GuildRequest struct {
	Page        int        `json:"p"`
	IDs         string     `json:"ids"`
	Name        string     `json:"name"`
	Established int        `json:"established"` // Year
	Since       int        `json:"since"`
	Status      string     `json:"status"`
	HasImages   string     `json:"hasImages"` // Y/N
	Order       GuildOrder `json:"order"`
	Sort        ListSort   `json:"sort"`
}

// GuildList represents a single "page" containing a slice of Guilds.
type GuildList struct {
	NumberOfPages int
	CurrentPage   int
	TotalResults  int
	Guilds        []Guild `json:"data"`
}

// List returns an GuildList containing a "page" of Guilds.
func (gs *GuildService) List(req *GuildRequest) (gl GuildList, err error) {
	// GET: /guilds
	v := encode(req)
	u := gs.c.url("/guilds", &v)

	var resp *http.Response
	resp, err = gs.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get guilds")
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&gl); err != nil {
		return
	}
	return
}

// Get retrieves a single Guild with the given guildID.
func (gs *GuildService) Get(guildID string) (g Guild, err error) {
	// GET: /guild/:guildID
	u := gs.c.url("/guild/"+guildID, nil)

	var resp *http.Response
	resp, err = gs.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get guild")
		return
	}
	defer resp.Body.Close()

	guildResponse := struct {
		Status  string
		Data    Guild
		Message string
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&guildResponse); err != nil {
		return
	}
	g = guildResponse.Data
	return
}

// AddGuild adds a Guild to the BreweryDB.
// The Guild Name is required.
func (gs *GuildService) AddGuild(g *Guild) (string, error) {
	// POST: /guilds
	v := encode(g)
	u := gs.c.url("/guilds", nil)

	resp, err := gs.c.PostForm(u, v)
	if err != nil {
		return "", err
	} else if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unable to add guild")
	}
	defer resp.Body.Close()

	addResponse := struct {
		Status string
		Data   struct {
			ID string
		}
		Message string
	}{}
	// TODO: return any response?
	if err := json.NewDecoder(resp.Body).Decode(&addResponse); err != nil {
		return "", err
	}

	return addResponse.Data.ID, nil
}

// UpdateGuild updates the Guild with the given guildID to match the given Guild.
func (gs *GuildService) UpdateGuild(guildID string, g *Guild) error {
	// PUT: /guild/:guildID
	u := gs.c.url("/guild/"+guildID, nil)
	v := encode(g)
	put, err := http.NewRequest("PUT", u, bytes.NewBufferString(v.Encode()))
	if err != nil {
		return err
	}

	resp, err := gs.c.Do(put)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to update guild")
	}
	defer resp.Body.Close()

	return nil
}

// DeleteGuild removes the Guild with the given guildID.
func (gs *GuildService) DeleteGuild(guildID string) error {
	// DELETE: /guild/:guildID
	u := gs.c.url("/guild/"+guildID, nil)

	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	resp, err := gs.c.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to delete guild")
	}
	defer resp.Body.Close()

	return nil
}
