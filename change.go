package brewerydb

import "net/http"

// ChangeService provides access to the BreweryDB Change API.
// Use Client.Change.
type ChangeService struct {
	c *Client
}

// ChangeType represents a type of attribute changed in BreweryDB.
type ChangeType string

// ChangeTypes available.
const (
	ChangeBeer     ChangeType = "beer"
	ChangeBrewery             = "brewery"
	ChangeEvent               = "event"
	ChangeGuild               = "guild"
	ChangeLocation            = "location"
)

// ChangeAction represents the type of change made to an attribute.
type ChangeAction string

// ChangeActions available.
const (
	ChangeAdd    ChangeAction = "add"
	ChangeDelete              = "delete"
	ChangeEdit                = "edit"
)

// ChangeListRequest contains parameters for obtaining a list of BreweryDB changes.
type ChangeListRequest struct {
	Page          int        `url:"p,omitempty"`
	AttributeName ChangeType `url:"attributeName,omitempty"`
	AttributeID   string     `url:"attributeId,omitempty"`
	Since         string     `url:"since,omitempty"`
}

// Attribute is a generic object that contains the ID and Name of either a
// Beer, Brewery, Event, Guild, or Location.
type Attribute struct {
	ID   string
	Name string
}

// Change contains all the relevant information for an individual change in BreweryDB.
type Change struct {
	AttributeName    ChangeType
	Action           ChangeAction
	Attribute        Attribute
	SubAttributeName ChangeType
	SubAction        ChangeAction
	SubAttribute     Attribute
}

// ChangeList represents one "page" containing a slice of Changes.
type ChangeList struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Changes       []Change `json:"data"`
}

// List retrieves (by default) a paginated list of all Changes to BreweryDB
// in the last 30 days.
func (cs *ChangeService) List(q *ChangeListRequest) (cl ChangeList, err error) {
	// GET: /changes
	var req *http.Request
	req, err = cs.c.NewRequest("GET", "/changes", q)
	if err != nil {
		return
	}

	err = cs.c.Do(req, &cl)
	return
}
