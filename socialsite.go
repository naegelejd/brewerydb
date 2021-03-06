package brewerydb

import (
	"fmt"
	"net/http"
)

// SocialSiteService provides access to the BreweryDB Social Site API.
// Use Client.SocialSite.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/socialsite_index
type SocialSiteService struct {
	c *Client
}

// SocialSite represents a social media website.
type SocialSite struct {
	ID         int
	Name       string
	Website    string
	CreateDate string
	UpdateDate string
}

// SocialAccount represents a social media account/handle.
// TODO: it appears some SocialAccount responses include the SocialSite ("socialMedia") object as well.
// TODO: SocialAccount responses also return an object corresponding to the query (e.g. Beer, Event, Guild, etc.)
type SocialAccount struct {
	ID            int        `url:"-"`
	SocialMediaID int        `url:"socialmediaId"`
	SocialSite    SocialSite `url:"-",json:"socialMedia"` // see TODO above
	Handle        string     `url:"handle"`
}

// List returns a slice of all SocialSites in the BreweryDB.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/socialsite_index#1
func (ss *SocialSiteService) List() (sl []SocialSite, err error) {
	// GET: /socialsites
	var req *http.Request
	req, err = ss.c.NewRequest("GET", "/socialsites", nil)
	if err != nil {
		return
	}

	socialsitesResp := struct {
		Status  string
		Data    []SocialSite
		Message string
	}{}
	err = ss.c.Do(req, &socialsitesResp)
	return socialsitesResp.Data, err
}

// Get retrieves the SocialSite having the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/socialsite_index#2
func (ss *SocialSiteService) Get(id int) (s SocialSite, err error) {
	// GET: /socialsite/:socialsiteId
	var req *http.Request
	req, err = ss.c.NewRequest("GET", fmt.Sprintf("/socialsite/%d", id), nil)
	if err != nil {
		return
	}

	socialsiteResp := struct {
		Status  string
		Data    SocialSite
		Message string
	}{}
	err = ss.c.Do(req, &socialsiteResp)
	return socialsiteResp.Data, err
}
