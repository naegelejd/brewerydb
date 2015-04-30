package brewerydb

import (
	"bytes"
	"encoding/json"
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

const (
	FestivalEvent            EventType = "festival"
	CompetitionEvent                   = "competition"
	FestivalCompetitionEvent           = "festival_competition"
	TastingEvent                       = "tasting"
	beerReleaseEvent                   = "beer_release"
	SeminarEvent                       = "seminar"
	MeetupEvent                        = "meetup"
	OtherEvent                         = "other"
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

const (
	WebsiteEventOrder        = "website"
	YearEventOrder           = "year"
	StartDateEventOrder      = "startDate"
	EndDateEventOrder        = "endDate"
	LocalityEventOrder       = "locality"
	RegionEventOrder         = "region"
	CountryIsoCodeEventOrder = "countryIsoCode"
	StatusEventOrder         = "status"
	CreateDateEventOrder     = "createDate"
	UpdateDateEventOrder     = "updateDate"
)

// EventList represents a single "page" containing a slice of Events.
type EventList struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Events        []Event `json:"data"`
}

// EventRequest contains options for specifying the kinds of Events desired.
type EventRequest struct {
	Page           int        `json:"p"`
	IDs            string     `json:"ids"`
	Year           string     `json:"year"`
	Name           string     `json:"name"`
	Type           string     `json:"type"`     // Key of the type of event, comma separated for multiple types
	Locality       string     `json:"locality"` // e.g. US city
	Region         string     `json:"region"`   // e.g. US state
	CountryISOCode string     `json:"countryIsoCode"`
	Since          int        `json:"since"` // Unix timestamp
	Status         string     `json:"status"`
	HasImages      string     `json:"hasImages"` // Y/N
	Order          EventOrder `json:"order"`
	Sort           ListSort   `json:"sort"`
}

// Events returns an EventList containing a "page" of Events.
func (es *EventService) Events(req *EventRequest) (el EventList, err error) {
	// GET: /events
	v := encode(req)
	u := es.c.url("/events", &v)

	var resp *http.Response
	resp, err = es.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get events")
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&el); err != nil {
		return
	}

	return
}

// Event retrieves a single event with the given eventID.
func (es *EventService) Event(eventID string) (e Event, err error) {
	// GET: /event/:eventID
	u := es.c.url("/event/"+eventID, nil)

	var resp *http.Response
	resp, err = es.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get event")
		return
	}
	defer resp.Body.Close()

	eventResponse := struct {
		Status  string
		Data    Event
		Message string
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&eventResponse); err != nil {
		return
	}
	e = eventResponse.Data

	return
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
	v := encode(e)
	u := es.c.url("/events", nil)

	resp, err := es.c.PostForm(u, v)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to add event")
	}
	defer resp.Body.Close()

	// TODO: return any response?

	return nil
}

// UpdateEvent updates the Event with the given eventID to match the given Event.
func (es *EventService) UpdateEvent(eventID string, e *Event) error {
	// PUT: /event/:eventID
	u := es.c.url("/event/"+eventID, nil)
	v := encode(e)
	put, err := http.NewRequest("PUT", u, bytes.NewBufferString(v.Encode()))
	if err != nil {
		return err
	}

	resp, err := es.c.Do(put)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to update event")
	}
	defer resp.Body.Close()

	// TODO: return any response?

	return nil
}

// DeleteEvent removes the Event with the given eventID.
func (es *EventService) DeleteEvent(eventID string) error {
	// DELETE: /event/:eventID
	u := es.c.url("/event/"+eventID, nil)

	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	resp, err := es.c.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to delete event")
	}
	defer resp.Body.Close()

	return nil
}
