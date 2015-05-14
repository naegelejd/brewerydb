package brewerydb

import "fmt"
import "net/http"

// TODO: according to http://www.brewerydb.com/developers/docs-endpoint/guild_socialaccount,
// the socialaccount requests return an object containing the SocialAccount fields, as well as
// an entire Guild object an entire SocialSite ("socialMedia") object!!!

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
