package brewerydb

import "net/http"

// ListBeers returns a page of Beers for the given Event.
func (es *EventService) ListBeers(eventID string) (bl BeerList, err error) {
	// GET: /event/:eventID/beers
	var req *http.Request
	req, err = es.c.NewRequest("GET", "/event/"+eventID+"/beers", nil)
	if err != nil {
		return
	}

	err = es.c.Do(req, &bl)
	return
}

// GetBeer retrieves the Beer with the given ID for the given Event.
func (es *EventService) GetBeer(eventID, beerID string) (b Beer, err error) {
	// GET: /event/:eventId/beer/:beerId
	var req *http.Request
	req, err = es.c.NewRequest("GET", "/event/"+eventID+"/beer/"+beerID, nil)
	if err != nil {
		return
	}

	eventBeerResp := struct {
		Status  string
		Data    Beer
		Message string
	}{}
	err = es.c.Do(req, &eventBeerResp)
	return eventBeerResp.Data, err
}

// EventChangeBeerRequest contains parameters for changing or adding
// a new Beer to an Event.
type EventChangeBeerRequest struct {
	IsPouring       string `json:"isPouring"`
	AwardCategoryID string `json:"awardCategoryId"`
	AwardPlaceID    string `json:"awardPlaceId"`
}

// AddBeer adds the Beer with the given ID to the given Event.
func (es *EventService) AddBeer(eventID, beerID string, q *EventChangeBeerRequest) error {
	// POST: /event/:eventId/beers
	params := struct {
		BeerID string `json:"beerId"`
		EventChangeBeerRequest
	}{beerID, *q}
	req, err := es.c.NewRequest("POST", "/event/"+eventID+"/beer/"+beerID, &params)
	if err != nil {
		return err
	}
	return es.c.Do(req, nil)
}

// UpdateBeer updates the Beer with the given ID for the given Event.
func (es *EventService) UpdateBeer(eventID, beerID string, q *EventChangeBeerRequest) error {
	// PUT: /event/:eventId/beer/:beerId
	req, err := es.c.NewRequest("PUT", "/event/"+eventID+"/beer/"+beerID, q)
	if err != nil {
		return err
	}
	return es.c.Do(req, nil)
}

// DeleteBeer removes the Beer with the given ID from the given Event.
func (es *EventService) DeleteBeer(eventID, beerID string) error {
	// DELETE: /event/:eventId/beer/:beerId
	req, err := es.c.NewRequest("DELETE", "/event/"+eventID+"/beer/"+beerID, nil)
	if err != nil {
		return err
	}
	return es.c.Do(req, nil)
}
