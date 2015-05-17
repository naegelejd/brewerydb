package brewerydb

import (
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
	ID          string `url:"-"`
	Name        string `url:"name"` // Required
	Description string `url:"description,omitempty"`
	Website     string `url:"website,omitempty"`
	Image       string `url:"image,omitempty"` // Base64. Only used for adding/updating Guilds.
	Images      struct {
		Icon   string `url:"-"`
		Medium string `url:"-"`
		Large  string `url:"-"`
	} `url:"-"`
	Established int `url:"established,omitempty"`
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

// GuildListRequest contains options for specifying kinds of Guilds desired.
type GuildListRequest struct {
	Page        int        `url:"p,omitempty"`
	IDs         string     `url:"ids,omitempty"`
	Name        string     `url:"name,omitempty"`        // Required for non-premium users.
	Established int        `url:"established,omitempty"` // Year
	Since       int        `url:"since,omitempty"`
	Status      string     `url:"status,omitempty"`
	HasImages   string     `url:"hasImages,omitempty"` // Y/N
	Order       GuildOrder `url:"order,omitempty"`
	Sort        ListSort   `url:"sort,omitempty"`
}

// GuildList represents a single "page" containing a slice of Guilds.
type GuildList struct {
	NumberOfPages int
	CurrentPage   int
	TotalResults  int
	Guilds        []Guild `json:"data"`
}

// List returns an GuildList containing a "page" of Guilds.
func (gs *GuildService) List(q *GuildListRequest) (gl GuildList, err error) {
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

// Add adds a Guild to the BreweryDB.
// The Guild Name is required.
func (gs *GuildService) Add(g *Guild) (string, error) {
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

// Update updates the Guild with the given guildID to match the given Guild.
func (gs *GuildService) Update(guildID string, g *Guild) error {
	// PUT: /guild/:guildID
	req, err := gs.c.NewRequest("PUT", "/guild/"+guildID, g)
	if err != nil {
		return err
	}

	return gs.c.Do(req, nil)
}

// Delete removes the Guild with the given guildID.
func (gs *GuildService) Delete(guildID string) error {
	// DELETE: /guild/:guildID
	req, err := gs.c.NewRequest("DELETE", "/guild/"+guildID, nil)
	if err != nil {
		return err
	}

	return gs.c.Do(req, nil)
}

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

// ListSocialAccounts returns a slice of all social media accounts associated with the given Guild.
func (gs *GuildService) ListSocialAccounts(guildID string) (sl []SocialAccount, err error) {
	// GET: /guild/:guildId/socialaccounts
	var req *http.Request
	req, err = gs.c.NewRequest("GET", "/guild/"+guildID+"/socialaccounts", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []SocialAccount
		Message string
	}{}
	err = gs.c.Do(req, &resp)
	return resp.Data, err
}

// GetSocialAccount retrieves the SocialAccount with the given ID for the given Guild.
func (gs *GuildService) GetSocialAccount(guildID string, socialAccountID int) (s SocialAccount, err error) {
	// GET: /guild/:guildId/socialaccount/:socialAccountId
	var req *http.Request
	req, err = gs.c.NewRequest("GET", fmt.Sprintf("/guild/%s/socialaccount/%d", guildID, socialAccountID), nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    SocialAccount
		Message string
	}{}
	err = gs.c.Do(req, &resp)
	return resp.Data, err
}

// AddSocialAccount adds a new SocialAccount to the given Guild.
func (gs *GuildService) AddSocialAccount(guildID string, s *SocialAccount) error {
	// POST: /guild/:guildId/socialaccounts
	req, err := gs.c.NewRequest("POST", "/guild/"+guildID+"/socialaccounts", s)
	if err != nil {
		return err
	}

	return gs.c.Do(req, nil)
}

// UpdateSocialAccount updates a SocialAccount for the given Guild.
func (gs *GuildService) UpdateSocialAccount(guildID string, s *SocialAccount) error {
	// PUT: /guild/:guildId/socialaccount/:socialAccountId
	req, err := gs.c.NewRequest("PUT", fmt.Sprintf("/guild/%s/socialaccount/%d", guildID, s.ID), s)
	if err != nil {
		return err
	}

	return gs.c.Do(req, nil)
}

// DeleteSocialAccount removes a SocialAccount from the given Guild.
func (gs *GuildService) DeleteSocialAccount(guildID string, socialAccountID int) error {
	// DELETE: /guild/:guildId/socialaccount/:socialAccountId
	req, err := gs.c.NewRequest("DELETE", fmt.Sprintf("/guild/%s/socialaccount/%d", guildID, socialAccountID), nil)
	if err != nil {
		return err
	}

	return gs.c.Do(req, nil)
}
