package brewerydb

import (
	"fmt"
	"net/http"
)

// ListSocialAccounts returns a slice of all social media accounts associated with the given Beer.
func (bs *BeerService) ListSocialAccounts(beerID string) (sl []SocialAccount, err error) {
	// GET: /beer/:beerId/socialaccounts
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+beerID+"/socialaccounts", nil)
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

// GetSocialAccount retrieves the SocialAccount with the given ID for the given Beer.
func (bs *BeerService) GetSocialAccount(beerID string, socialAccountID int) (s SocialAccount, err error) {
	// GET: /beer/:beerId/socialaccount/:socialaccountId
	var req *http.Request
	req, err = bs.c.NewRequest("GET", fmt.Sprintf("/beer/%s/socialaccount/%d", beerID, socialAccountID), nil)
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

// AddSocialAccount adds a new SocialAccount to the given Beer.
func (bs *BeerService) AddSocialAccount(beerID string, s *SocialAccount) error {
	// POST: /beer/:beerId/socialaccounts
	req, err := bs.c.NewRequest("POST", "/beer/"+beerID+"/socialaccounts", s)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// UpdateSocialAccount updates a SocialAccount for the given Beer.
func (bs *BeerService) UpdateSocialAccount(beerID string, s *SocialAccount) error {
	// PUT: /beer/:beerId/socialaccount/:socialaccountId
	req, err := bs.c.NewRequest("PUT", fmt.Sprintf("/beer/%s/socialaccount/%d", beerID, s.ID), s)
	if err != nil {
		return err
	}

	return bs.c.Do(req, nil)
}

// DeleteSocialAccount removes a SocialAccount from the given Beer.
func (bs *BeerService) DeleteSocialAccount(beerID string, socialAccountID int) error {
	// DELETE: /beer/:beerId/socialaccount/:socialaccountId
	req, err := bs.c.NewRequest("DELETE", fmt.Sprintf("/beer/%s/socialaccount/%d", beerID, socialAccountID), nil)
	if err != nil {
		return err
	}
	return bs.c.Do(req, nil)
}
