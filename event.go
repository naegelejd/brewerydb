package brewerydb

import "net/http"

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
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Type            EventType `json:"type"`
	StartDate       string    `json:"startDate"` // YYYY-MM-DD
	EndDate         string    `json:"endDate"`   // YYYY-MM-DD
	Description     string    `json:"description"`
	Year            string    `json:"year"`
	Time            string    `json:"time"`
	Price           string    `json:"price"`
	VenueName       string    `json:"venueName"`
	StreetAddress   string    `json:"streetAddress"`
	ExtendedAddress string    `json:"extendedAddress"`
	Locality        string    `json:"locality"`
	Region          string    `json:"region"`
	PostalCode      string    `json:"postalCode"`
	CountryIsoCode  string    `json:"countryIsoCode"` // Required
	Phone           string    `json:"phone"`
	Website         string    `json:"website"`
	Longitude       float64   `json:"longitude"`
	Latitude        float64   `json:"latitude"`
	Image           string    `json:"image"` // base64-encoded

	// TODO: The following should be empty when adding or updating an Event
	Images struct {
		Icon   string `json:",omitempty"`
		Medium string `json:",omitempty"`
		Large  string `json:",omitempty"`
	} `json:",omitempty"`
	CreateDate    string `json:",omitempty"`
	UpdateDate    string `json:",omitempty"`
	Status        string `json:",omitempty"`
	StatusDisplay string `json:",omitempty"`
	Country       struct {
		Name        string `json:",omitempty"`
		DisplayName string `json:",omitempty"`
		ISOCode     string `json:",omitempty"`
		ISOThree    string `json:",omitempty"`
		NumberCode  int    `json:",omitempty"`
		CreateDate  string `json:",omitempty"`
	} `json:",omitempty"`
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

// EventRequest contains options for specifying the kinds of Events desired.
// Non-Premium users must set one of the following: year, name, type, locality, region
type EventRequest struct {
	Page           int        `json:"p"`
	IDs            string     `json:"ids,omitempty"`
	Year           int        `json:"year,omitempty"`
	Name           string     `json:"name,omitempty"`
	Type           string     `json:"type,omitempty"`     // Key of the type of event, comma separated for multiple types
	Locality       string     `json:"locality,omitempty"` // e.g. US city
	Region         string     `json:"region,omitempty"`   // e.g. US state
	CountryISOCode string     `json:"countryIsoCode,omitempty"`
	Since          int        `json:"since,omitempty"` // Unix timestamp
	Status         string     `json:"status,omitempty"`
	HasImages      string     `json:"hasImages,omitempty"` // Y/N
	Order          EventOrder `json:"order,omitempty"`
	Sort           ListSort   `json:"sort,omitempty"`
}

// List returns an EventList containing a "page" of Events.
func (es *EventService) List(q *EventRequest) (el EventList, err error) {
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

	eventResponse := struct {
		Status  string
		Data    Event
		Message string
	}{}

	if err = es.c.Do(req, &eventResponse); err != nil {
		return
	}

	return eventResponse.Data, nil
}

// AddEvent adds an Event to the BreweryDB.
// The following **must** be set in the Event:
//
// - Name
// - Type
// - StartDate (YYYY-MM-DD)
// - EndDate (YYYY-MM-DD)
func (es *EventService) AddEvent(e *Event) error {
	// POST: /events
	req, err := es.c.NewRequest("POST", "/events", e)
	if err != nil {
		return err
	}

	// TODO: return any response?
	return es.c.Do(req, nil)
}

// UpdateEvent updates the Event with the given eventID to match the given Event.
func (es *EventService) UpdateEvent(eventID string, e *Event) error {
	// PUT: /event/:eventID
	req, err := es.c.NewRequest("PUT", "/event/"+eventID, e)
	if err != nil {
		return err
	}

	// TODO: return any response?
	return es.c.Do(req, nil)
}

// DeleteEvent removes the Event with the given eventID.
func (es *EventService) DeleteEvent(eventID string) error {
	// DELETE: /event/:eventID
	req, err := es.c.NewRequest("DELETE", "/event/"+eventID, nil)
	if err != nil {
		return err
	}

	return es.c.Do(req, nil)
}
