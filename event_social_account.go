package brewerydb

import (
	"fmt"
	"net/http"
)

// SocialAccount represents a social media account/handle.
type SocialAccount struct {
	ID            int
	SocialMediaID int    `json:"socialmediaId"`
	Handle        string `json:"handle"`
}

// ListSocialAccounts returns a slice of all social media accounts associated with the given Event.
func (es *EventService) ListSocialAccounts(eventID string) (sl []SocialAccount, err error) {
	// GET: /event/:eventId/socialaccounts
	var req *http.Request
	req, err = es.c.NewRequest("GET", "/event/"+eventID+"/socialaccounts", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []SocialAccount
		Message string
	}{}
	err = es.c.Do(req, &resp)
	return resp.Data, err
}

// GetSocialAccount retrieves the SocialAccount with the given ID for the given Event.
func (es *EventService) GetSocialAccount(eventID string, socialAccountID int) (s SocialAccount, err error) {
	// GET: /event/:eventId/socialaccount/:socialaccountId
	var req *http.Request
	req, err = es.c.NewRequest("GET", fmt.Sprintf("/event/%s/socialaccount/%d", eventID, socialAccountID), nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    SocialAccount
		Message string
	}{}
	err = es.c.Do(req, &resp)
	return resp.Data, err
}

// AddSocialAccount adds a new SocialAccount to the given Event.
// TODO: this encodes the entire SocialAccount, but its
// ID should not be url-encoded
func (es *EventService) AddSocialAccount(eventID string, s *SocialAccount) error {
	// POST: /event/:eventId/socialaccounts
	req, err := es.c.NewRequest("POST", "/event/"+eventID+"/socialaccounts", s)
	if err != nil {
		return err
	}

	return es.c.Do(req, nil)
}

// UpdateSocialAccount updates a SocialAccount for the given Event.
// TODO: this encodes the entire SocialAccount, but its
// ID should not be url-encoded
func (es *EventService) UpdateSocialAccount(eventID string, s *SocialAccount) error {
	// PUT: /event/:eventId/socialaccount/:socialaccountId
	req, err := es.c.NewRequest("PUT", fmt.Sprintf("/event/%s/socialaccount/%d", eventID, s.ID), s)
	if err != nil {
		return err
	}

	return es.c.Do(req, nil)
}

// DeleteSocialAccount removes a SocialAccount from the given Event.
func (es *EventService) DeleteSocialAccount(eventID string, socialAccountID int) error {
	// DELETE: /event/:eventId/socialaccount/:socialaccountId
	req, err := es.c.NewRequest("DELETE", fmt.Sprintf("/event/%s/socialaccount/%d", eventID, socialAccountID), nil)
	if err != nil {
		return err
	}
	return es.c.Do(req, nil)
}
