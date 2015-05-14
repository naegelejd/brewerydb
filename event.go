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
	CountryIsoCode  string    `url:"countryIsoCode"` // Required
	Phone           string    `url:"phone,omitempty"`
	Website         string    `url:"website,omitempty"`
	Longitude       float64   `url:"longitude,omitempty"`
	Latitude        float64   `url:"latitude,omitempty"`
	Image           string    `url:"image"` // base64. Only used for adding/updating Events.
	Images          struct {
		Icon   string `url:"-"`
		Medium string `url:"-"`
		Large  string `url:"-"`
	} `url:"-"`
	Status        string  `url:"-"`
	StatusDisplay string  `url:"-"`
	Country       Country `url:"-"`
	CreateDate    string  `url:"-"`
	UpdateDate    string  `url:"-"`
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
	HasImages      string     `url:"hasImages,omitempty"` // Y/N
	Order          EventOrder `url:"order,omitempty"`
	Sort           ListSort   `url:"sort,omitempty"`
}

// List returns an EventList containing a "page" of Events.
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
