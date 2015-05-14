package brewerydb

import (
	"fmt"
	"net/http"
)

// ListSocialAccounts returns a slice of all social media accounts associated with the given Brewery.
func (bs *BreweryService) ListSocialAccounts(breweryID string) (sl []SocialAccount, err error) {
	// GET: /brewery/:breweryId/socialaccounts
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/brewery/"+breweryID+"/socialaccounts", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []SocialAccount
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// GetSocialAccount retrieves the SocialAccount with the given ID for the given Brewery.
func (bs *BreweryService) GetSocialAccount(breweryID string, socialAccountID int) (s SocialAccount, err error) {
	// GET: /brewery/:breweryId/socialaccount/:socialaccountId
	var req *http.Request
	req, err = bs.c.NewRequest("GET", fmt.Sprintf("/brewery/%s/socialaccount/%d", breweryID, socialAccountID), nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    SocialAccount
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}

// AddSocialAccount adds a new SocialAccount to the given Brewery.
func (bs *BreweryService) AddSocialAccount(breweryID string, s *SocialAccount) error {
	// POST: /brewery/:breweryId/socialaccounts
	req, err := bs.c.NewRequest("POST", "/brewery/"+breweryID+"/socialaccounts", s)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// UpdateSocialAccount updates a SocialAccount for the given Brewery.
func (bs *BreweryService) UpdateSocialAccount(breweryID string, s *SocialAccount) error {
	// PUT: /brewery/:breweryId/socialaccount/:socialaccountId
	req, err := bs.c.NewRequest("PUT", fmt.Sprintf("/brewery/%s/socialaccount/%d", breweryID, s.ID), s)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// DeleteSocialAccount removes a SocialAccount from the given Brewery.
func (bs *BreweryService) DeleteSocialAccount(breweryID string, socialAccountID int) error {
	// DELETE: /brewery/:breweryId/socialaccount/:socialaccountId
	req, err := bs.c.NewRequest("DELETE", fmt.Sprintf("/brewery/%s/socialaccount/%d", breweryID, socialAccountID), nil)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}
