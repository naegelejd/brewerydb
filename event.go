package brewerydb

import (
	"fmt"
	"net/http"
)

// EventService provides access to the BreweryDB Event API.
// Use Client.Event.
type EventService struct {
	c *Client
}

// EventType specifies the type of the Event.
type EventType string

// Types of Events.
const (
	EventFestival            EventType = "festival"
	EventCompetition                   = "competition"
	EventFestivalCompetition           = "festival_competition"
	EventTasting                       = "tasting"
	EventBeerRelease                   = "beer_release"
	EventSeminar                       = "seminar"
	EventMeetup                        = "meetup"
	EventOther                         = "other"
)

// Event represents a community event related to Beer/Breweries.
type Event struct {
	ID              string    `url:"-"`
	Name            string    `url:"name"`
	Type            EventType `url:"type"`
	StartDate       string    `url:"startDate"` // YYYY-MM-DD
	EndDate         string    `url:"endDate"`   // YYYY-MM-DD
	Description     string    `url:"description,omitempty"`
	Year            string    `url:"year,omitempty"`
	Time            string    `url:"time,omitempty"`
	Price           string    `url:"price,omitempty"`
	VenueName       string    `url:"venueName,omitempty"`
	StreetAddress   string    `url:"streetAddress,omitempty"`
	ExtendedAddress string    `url:"extendedAddress,omitempty"`
	Locality        string    `url:"locality,omitempty"`
	Region          string    `url:"region,omitempty"`
	PostalCode      string    `url:"postalCode,omitempty"`
	CountryISOCode  string    `url:"countryIsoCode"` // Required
	Phone           string    `url:"phone,omitempty"`
	Website         string    `url:"website,omitempty"`
	Longitude       float64   `url:"longitude,omitempty"`
	Latitude        float64   `url:"latitude,omitempty"`
	Image           string    `url:"image"` // base64. Only used for adding/updating Events.
	Images          Images    `url:"-"`
	Status          string    `url:"-"`
	StatusDisplay   string    `url:"-"`
	Country         Country   `url:"-"`
	CreateDate      string    `url:"-"`
	UpdateDate      string    `url:"-"`
}

// EventOrder specifies the ordering of an EventList.
type EventOrder string

// EventList ordering options.
const (
	EventOrderWebsite        EventOrder = "website"
	EventOrderYear                      = "year"
	EventOrderStartDate                 = "startDate"
	EventOrderEndDate                   = "endDate"
	EventOrderLocality                  = "locality"
	EventOrderRegion                    = "region"
	EventOrderCountryIsoCode            = "countryIsoCode"
	EventOrderStatus                    = "status"
	EventOrderCreateDate                = "createDate"
	EventOrderUpdateDate                = "updateDate"
)

// EventList represents a single "page" containing a slice of Events.
type EventList struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Events        []Event `json:"data"`
}

// EventListRequest contains options for specifying the kinds of Events desired.
// Non-Premium users must set one of the following: year, name, type, locality, region
type EventListRequest struct {
	Page           int        `url:"p"`
	IDs            string     `url:"ids,omitempty"`
	Year           int        `url:"year,omitempty"`
	Name           string     `url:"name,omitempty"`
	Type           string     `url:"type,omitempty"`     // Key of the type of event, comma separated for multiple types
	Locality       string     `url:"locality,omitempty"` // e.g. US city
	Region         string     `url:"region,omitempty"`   // e.g. US state
	CountryISOCode string     `url:"countryIsoCode,omitempty"`
	Since          int        `url:"since,omitempty"` // Unix timestamp
	Status         string     `url:"status,omitempty"`
	HasImages      YesNo      `url:"hasImages,omitempty"`
	Order          EventOrder `url:"order,omitempty"`
	Sort           ListSort   `url:"sort,omitempty"`
}

// List returns an EventList containing a "page" of Events.
// For non-premium members, one of Year, Name, Type, Locality or Region must be set.
func (es *EventService) List(q *EventListRequest) (el EventList, err error) {
	// GET: /events
	var req *http.Request
	req, err = es.c.NewRequest("GET", "/events", q)
	if err != nil {
		return
	}

	err = es.c.Do(req, &el)
	return
}

// Get retrieves a single event with the given eventID.
func (es *EventService) Get(eventID string) (e Event, err error) {
	// GET: /event/:eventID
	var req *http.Request
	req, err = es.c.NewRequest("GET", "/event/"+eventID, nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    Event
		Message string
	}{}

	err = es.c.Do(req, &resp)
	return resp.Data, err
}

// Add adds an Event to the BreweryDB and returns its new ID.
// The following **must** be set in the Event:
//
// - Name
// - Type
// - StartDate (YYYY-MM-DD)
// - EndDate (YYYY-MM-DD)
func (es *EventService) Add(e *Event) (string, error) {
	// POST: /events
	if e == nil {
		return "", fmt.Errorf("nil Event")
	}
	req, err := es.c.NewRequest("POST", "/events", e)
	if err != nil {
		return "", err
	}

	resp := struct {
		Status string
		Data   struct {
			ID string
		}
		Message string
	}{}

	err = es.c.Do(req, &resp)
	return resp.Data.ID, err
}

// Update updates the Event with the given eventID to match the given Event.
func (es *EventService) Update(eventID string, e *Event) error {
	// PUT: /event/:eventID
	if e == nil {
		return fmt.Errorf("nil Event")
	}
	req, err := es.c.NewRequest("PUT", "/event/"+eventID, e)
	if err != nil {
		return err
	}

	// TODO: return any response?
	return es.c.Do(req, nil)
}

// Delete removes the Event with the given eventID.
func (es *EventService) Delete(eventID string) error {
	// DELETE: /event/:eventID
	req, err := es.c.NewRequest("DELETE", "/event/"+eventID, nil)
	if err != nil {
		return err
	}

	return es.c.Do(req, nil)
}

// AwardCategory represents a category of award for an Event.
type AwardCategory struct {
	ID          int    `url:"-"`
	Name        string `url:"name"` // required for adding/updating AwardCategories
	Description string `url:"description"`
	Image       string `url:"image"` // base64
	CreateDate  string `url:"-"`
	UpdateDate  string `url:"-"`
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
	if a == nil {
		return fmt.Errorf("nil AwardCategory")
	}
	req, err := es.c.NewRequest("POST", "/event/"+eventID+"/awardcategories", a)
	if err != nil {
		return err
	}

	return es.c.Do(req, nil)
}

// UpdateAwardCategory updates an AwardCategory	for the given Event.
func (es *EventService) UpdateAwardCategory(eventID string, a *AwardCategory) error {
	// PUT: /event/:eventId/awardcategory/:awardcategoryId
	if a == nil {
		return fmt.Errorf("nil AwardCategory")
	}
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

// AwardPlace represents an award location.
type AwardPlace struct {
	ID          int    `url:"-"`
	Name        string `url:"name"` // required for adding/updating AwardPlaces
	Description string `url:"description"`
	Image       string `url:"image"` // base64
	CreateDate  string `url:"-"`
	UpdateDate  string `url:"-"`
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
	if a == nil {
		return fmt.Errorf("nil AwardPlace")
	}
	req, err := es.c.NewRequest("POST", "/event/"+eventID+"/awardplaces", a)
	if err != nil {
		return err
	}

	return es.c.Do(req, nil)
}

// UpdateAwardPlace updates an AwardPlace for the given Event.
func (es *EventService) UpdateAwardPlace(eventID string, a *AwardPlace) error {
	// PUT: /event/:eventId/awardplace/:awardplaceId
	if a == nil {
		return fmt.Errorf("nil AwardPlace")
	}
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

// EventBeersRequest contains parameters for specifying desired
// Beers for a given Event.
type EventBeersRequest struct {
	Page            int   `url:"p, omitempty"`
	OnlyWinners     YesNo `url:"onlyWinners,omitempty"`
	AwardCategoryID int   `url:"awardcategoryId,omitempty"`
	AwardPlaceID    int   `url:"awardplaceId,omitempty"`
}

// ListBeers returns a page of Beers for the given Event.
func (es *EventService) ListBeers(eventID string, q *EventBeersRequest) (bl BeerList, err error) {
	// GET: /event/:eventID/beers
	var req *http.Request
	req, err = es.c.NewRequest("GET", "/event/"+eventID+"/beers", q)
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
	IsPouring       YesNo `url:"isPouring,omitempty"`
	AwardCategoryID int   `url:"awardcategoryId,omitempty"`
	AwardPlaceID    int   `url:"awardplaceId,omitempty"`
}

type eventAddBeerRequest struct {
	BeerID string `url:"beerId"`
	EventChangeBeerRequest
}

// AddBeer adds the Beer with the given ID to the given Event.
func (es *EventService) AddBeer(eventID, beerID string, q *EventChangeBeerRequest) error {
	// POST: /event/:eventId/beers
	var params *eventAddBeerRequest
	if q != nil {
		params = &eventAddBeerRequest{beerID, *q}
	} else {
		params = &eventAddBeerRequest{BeerID: beerID}
	}
	req, err := es.c.NewRequest("POST", "/event/"+eventID+"/beers", params)
	if err != nil {
		return err
	}
	return es.c.Do(req, nil)
}

// UpdateBeer updates the Beer with the given ID for the given Event.
func (es *EventService) UpdateBeer(eventID, beerID string, q *EventChangeBeerRequest) error {
	// PUT: /event/:eventId/beer/:beerId
	if q == nil {
		q = &EventChangeBeerRequest{}
	}
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

// EventBreweriesRequest contains parameters for specifying desired
// Breweries for a given Event.
type EventBreweriesRequest struct {
	Page            int   `url:"p,omitempty"`
	OnlyWinners     YesNo `url:"onlyWinners,omitempty"`
	AwardCategoryID int   `url:"awardcategoryId,omitempty"`
	AwardPlaceID    int   `url:"awardplaceId,omitempty"`
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
	AwardCategoryID int `url:"awardcategoryId,omitempty"`
	AwardPlaceID    int `url:"awardplaceId,omitempty"`
}

// TODO: test the encoding of embedded structs:
type eventAddBreweryRequest struct {
	BreweryID string `url:"breweryId"`
	EventChangeBreweryRequest
}

// AddBrewery adds the Brewery with the given ID to the given Event.
func (es *EventService) AddBrewery(eventID, breweryID string, q *EventChangeBreweryRequest) error {
	// POST: /event/:eventID/brewery/:breweryID
	var params *eventAddBreweryRequest
	if q != nil {
		params = &eventAddBreweryRequest{breweryID, *q}
	} else {
		params = &eventAddBreweryRequest{BreweryID: breweryID}
	}

	req, err := es.c.NewRequest("POST", "/event/"+eventID+"/breweries", params)
	if err != nil {
		return err
	}
	return es.c.Do(req, nil)
}

// UpdateBrewery updates the Brewery with the given ID for the given Event.
func (es *EventService) UpdateBrewery(eventID, breweryID string, q *EventChangeBreweryRequest) error {
	// PUT: /event/:eventID/brewery/:breweryID
	if q == nil {
		q = &EventChangeBreweryRequest{}
	}
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
func (es *EventService) AddSocialAccount(eventID string, s *SocialAccount) error {
	// POST: /event/:eventId/socialaccounts
	if s == nil {
		return fmt.Errorf("nil SocialAccount")
	}
	req, err := es.c.NewRequest("POST", "/event/"+eventID+"/socialaccounts", s)
	if err != nil {
		return err
	}

	return es.c.Do(req, nil)
}

// UpdateSocialAccount updates a SocialAccount for the given Event.
func (es *EventService) UpdateSocialAccount(eventID string, s *SocialAccount) error {
	// PUT: /event/:eventId/socialaccount/:socialaccountId
	if s == nil {
		return fmt.Errorf("nil SocialAccount")
	}
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
