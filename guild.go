package brewerydb

import "net/http"

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
func (gs *GuildService) List(q *GuildRequest) (gl GuildList, err error) {
	// GET: /guilds
	var req *http.Request
	req, err = gs.c.NewRequest("GET", "/guilds", q)
	if err != nil {
		return
	}

	err = gs.c.Do(req, &gl)
	return
}

// Get retrieves a single Guild with the given guildID.
func (gs *GuildService) Get(guildID string) (g Guild, err error) {
	// GET: /guild/:guildID
	var req *http.Request
	req, err = gs.c.NewRequest("GET", "/guild/"+guildID, nil)
	if err != nil {
		return
	}

	guildResponse := struct {
		Status  string
		Data    Guild
		Message string
	}{}
	err = gs.c.Do(req, &guildResponse)
	return guildResponse.Data, err
}

// AddGuild adds a Guild to the BreweryDB.
// The Guild Name is required.
func (gs *GuildService) AddGuild(g *Guild) (string, error) {
	// POST: /guilds
	req, err := gs.c.NewRequest("POST", "/guilds", g)
	if err != nil {
		return "", err
	}

	addResponse := struct {
		Status string
		Data   struct {
			ID string
		}
		Message string
	}{}
	// TODO: return any response?
	err = gs.c.Do(req, &addResponse)
	return addResponse.Data.ID, err
}

// UpdateGuild updates the Guild with the given guildID to match the given Guild.
func (gs *GuildService) UpdateGuild(guildID string, g *Guild) error {
	// PUT: /guild/:guildID
	req, err := gs.c.NewRequest("PUT", "/guild/"+guildID, g)
	if err != nil {
		return err
	}

	return gs.c.Do(req, nil)
}

// DeleteGuild removes the Guild with the given guildID.
func (gs *GuildService) DeleteGuild(guildID string) error {
	// DELETE: /guild/:guildID
	req, err := gs.c.NewRequest("DELETE", "/guild/"+guildID, nil)
	if err != nil {
		return err
	}

	return gs.c.Do(req, nil)
}
