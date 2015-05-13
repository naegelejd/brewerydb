package brewerydb

import "net/http"

// EventBreweriesRequest contains parameters for specifying desired
// Breweries for a given Event.
type EventBreweriesRequest struct {
	Page            int    `json:"p"`
	OnlyWinners     string `json:"onlyWinners,omitempty"` // Y/N
	AwardCategoryID string `json:"awardCategoryId,omitempty"`
	AwardPlaceID    string `json:"awardPlaceId,omitempty"`
}

// ListBreweries returns a page of Breweries for the given Event.
func (es *EventService) ListBreweries(eventID string, q *EventBreweriesRequest) (bl BreweryList, err error) {
	// GET: /event/:eventID/breweries
	var req *http.Request
	req, err = es.c.NewRequest("GET", "/event/"+eventID+"/breweries", q)
	if err != nil {
		return
	}

	err = es.c.Do(req, &bl)
	return
}

// GetBrewery retrieves the Brewery with the given ID for the given Event.
func (es *EventService) GetBrewery(eventID, breweryID string) (b Brewery, err error) {
	// GET: /event/:eventID/brewery/:breweryID
	var req *http.Request
	req, err = es.c.NewRequest("GET", "/event/"+eventID+"/brewery/"+breweryID, nil)
	if err != nil {
		return
	}

	eventBreweryResp := struct {
		Status  string
		Data    Brewery
		Message string
	}{}
	err = es.c.Do(req, &eventBreweryResp)
	return eventBreweryResp.Data, err
}

// EventChangeBreweryRequest contains parameters for changing or adding
// a new Brewery to an Event.
type EventChangeBreweryRequest struct {
	AwardCategoryID string `json:"awardCategoryId,omitempty"`
	AwardPlaceID    string `json:"awardPlaceId,omitempty"`
}

// AddBrewery adds the Brewery with the given ID to the given Event.
func (es *EventService) AddBrewery(eventID, breweryID string, q *EventChangeBreweryRequest) error {
	// POST: /event/:eventID/brewery/:breweryID
	// TODO: test the encoding of embedded structs:
	params := struct {
		BreweryID string `json:"breweryId"`
		EventChangeBreweryRequest
	}{breweryID, *q}
	req, err := es.c.NewRequest("POST", "/event/"+eventID+"/brewery/"+breweryID, &params)
	if err != nil {
		return err
	}
	return es.c.Do(req, nil)
}

// UpdateBrewery updates the Brewery with the given ID for the given Event.
func (es *EventService) UpdateBrewery(eventID, breweryID string, q *EventChangeBreweryRequest) error {
	// PUT: /event/:eventID/brewery/:breweryID
	req, err := es.c.NewRequest("PUT", "/event/"+eventID+"/brewery/"+breweryID, q)
	if err != nil {
		return err
	}
	return es.c.Do(req, nil)
}

// DeleteBrewery removes the Brewery with the given ID from the given Event.
func (es *EventService) DeleteBrewery(eventID, breweryID string) error {
	// DELETE: /event/:eventID/brewery/:breweryID
	req, err := es.c.NewRequest("DELETE", "/event/"+eventID+"/brewery/"+breweryID, nil)
	if err != nil {
		return err
	}
	return es.c.Do(req, nil)
}
