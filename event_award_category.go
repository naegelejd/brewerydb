package brewerydb

import (
	"fmt"
	"net/http"
)

// AwardCategory represents a category of award for an Event.
type AwardCategory struct {
	ID          int
	Name        string
	Description string
	Image       string // base64
	CreateDate  string
	UpdateDate  string
}

// ListAwardCategories returns a slice of all AwardCategories for the given Event.
func (es *EventService) ListAwardCategories(eventID string) (al []AwardCategory, err error) {
	// GET: /event/:eventId/awardcategories
	var req *http.Request
	req, err = es.c.NewRequest("GET", "/event/"+eventID+"/awardcategories", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []AwardCategory
		Message string
	}{}
	err = es.c.Do(req, &resp)
	return resp.Data, nil
}

// GetAwardCategory retrieves the specified AwardCategory for the given Event.
func (es *EventService) GetAwardCategory(eventID string, awardCategoryID int) (a AwardCategory, err error) {
	// GET: /event/:eventId/awardcategory/:awardcategoryId
	var req *http.Request
	req, err = es.c.NewRequest("GET", fmt.Sprintf("/event/%s/awardcategory/%d", eventID, awardCategoryID), nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    AwardCategory
		Message string
	}{}
	err = es.c.Do(req, &resp)
	return resp.Data, err
}

// AddAwardCategory adds a new AwardCategory to the given Event.
func (es *EventService) AddAwardCategory(eventID string, a *AwardCategory) error {
	// POST: /event/:eventId/awardcategories
	req, err := es.c.NewRequest("POST", "/event/"+eventID+"/awardcategories", a)
	if err != nil {
		return err
	}

	return es.c.Do(req, nil)
}

// UpdateAwardCategory updates an AwardCategory	for the given Event.
func (es *EventService) UpdateAwardCategory(eventID string, a *AwardCategory) error {
	// PUT: /event/:eventId/awardcategory/:awardcategoryId
	req, err := es.c.NewRequest("PUT", fmt.Sprintf("/event/%s/awardcategory/%d", eventID, a.ID), a)
	if err != nil {
		return err
	}

	return es.c.Do(req, nil)
}

// DeleteAwardCategory removes an AwardCategory from the given Event.
func (es *EventService) DeleteAwardCategory(eventID string, awardCategoryID int) error {
	// DELETE: /event/:eventId/awardcategory/:awardcategoryId
	req, err := es.c.NewRequest("DELETE", fmt.Sprintf("/event/%s/awardcategory/%d", eventID, awardCategoryID), nil)
	if err != nil {
		return err
	}

	return es.c.Do(req, nil)
}
