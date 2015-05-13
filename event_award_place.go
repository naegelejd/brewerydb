package brewerydb

import (
	"fmt"
	"net/http"
)

// AwardPlace represents an award location.
type AwardPlace struct {
	ID          int
	Name        string
	Description string
	Image       string // base64
	CreateDate  string
	UpdateDate  string
}

// ListAwardPlaces returns a slice of all AwardPlaces for the given Event.
func (es *EventService) ListAwardPlaces(eventID string) (al []AwardPlace, err error) {
	// GET: /event/:eventId/awardplaces
	var req *http.Request
	req, err = es.c.NewRequest("GET", "/event/"+eventID+"/awardplaces", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []AwardPlace
		Message string
	}{}
	err = es.c.Do(req, &resp)
	return resp.Data, nil
}

// GetAwardPlace retrieves the specified AwardPlace for the given Event.
func (es *EventService) GetAwardPlace(eventID string, awardPlaceID int) (a AwardPlace, err error) {
	// GET: /event/:eventId/awardplace/:awardplaceId
	var req *http.Request
	req, err = es.c.NewRequest("GET", fmt.Sprintf("/event/%s/awardplace/%d", eventID, awardPlaceID), nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    AwardPlace
		Message string
	}{}
	err = es.c.Do(req, &resp)
	return resp.Data, err
}

// AddAwardPlace adds a new AwardPlace to the given Event.
func (es *EventService) AddAwardPlace(eventID string, a *AwardPlace) error {
	// POST: /event/:eventId/awardplaces
	req, err := es.c.NewRequest("POST", "/event/"+eventID+"/awardplaces", a)
	if err != nil {
		return err
	}

	return es.c.Do(req, nil)
}

// UpdateAwardPlace updates an AwardPlace for the given Event.
func (es *EventService) UpdateAwardPlace(eventID string, a *AwardPlace) error {
	// PUT: /event/:eventId/awardplace/:awardplaceId
	req, err := es.c.NewRequest("PUT", fmt.Sprintf("/event/%s/awardplace/%d", eventID, a.ID), a)
	if err != nil {
		return err
	}

	return es.c.Do(req, nil)
}

// DeleteAwardPlace removes an AwardPlace from the given Event.
func (es *EventService) DeleteAwardPlace(eventID string, awardPlaceID int) error {
	// DELETE: /event/:eventId/awardplace/:awardplaceId
	req, err := es.c.NewRequest("DELETE", fmt.Sprintf("/event/%s/awardplace/%d", eventID, awardPlaceID), nil)
	if err != nil {
		return err
	}

	return es.c.Do(req, nil)
}
