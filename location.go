package brewerydb

import "fmt"
import "net/http"

// LocationService provides access to the BreweryDB Location API.
// Use Client.Location.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/location_index
type LocationService struct {
	c *Client
}

// Country represents a country found on Earth.
type Country struct {
	ISOCode     string
	Name        string
	DisplayName string
	ISOThree    string
	NumberCode  int
	URLTitle    string
	CreateDate  string
	// UpdateDate  string
}

// Location represents a the location of a Brewery or similar institution.
type Location struct {
	ID                       string       `url:"-"`
	Name                     string       `url:"name,omitempty"`
	StreetAddress            string       `url:"streetAddress,omitempty"`
	ExtendedAddress          string       `url:"extendedAddress,omitempty"`
	Locality                 string       `url:"locality,omitempty"`
	Region                   string       `url:"region,omitempty"`
	PostalCode               string       `url:"postalCode,omitempty"`
	Phone                    string       `url:"phone,omitempty"`
	Website                  string       `url:"website,omitempty"`
	HoursOfOperation         string       `url:"-"`
	HoursOfOperationExplicit []string     `url:"hoursOfOperationExplicit,omitempty"`
	HoursOfOperationNotes    string       `url:"hoursOfOperationNotes,omitempty"`
	TourInfo                 string       `url:"tourInfo,omitempty"`
	TimezoneID               string       `url:"timezoneId,omitempty"`
	Latitude                 float64      `url:"latitude,omitempty"`
	Longitude                float64      `url:"longitude,omitempty"`
	IsPrimary                YesNo        `url:"isPrimary,omitempty"`
	InPlanning               YesNo        `url:"inPlanning,omitempty"`
	IsClosed                 YesNo        `url:"isClosed,omitempty"`
	OpenToPublic             YesNo        `url:"openToPublic,omitempty"`
	LocationType             LocationType `url:"locationType,omitempty"`
	LocationTypeDisplay      string       `url:"-"`
	CountryISOCode           string       `url:"countryIsoCode"` // Required for UpdateLocation
	Country                  Country      `url:"-"`
	CreateDate               string       `url:"-"`
	UpdateDate               string       `url:"-"`
	YearOpened               string       `url:"-"`
	YearClosed               string       `url:"-"`
	BreweryID                string       `url:"-"`
	Brewery                  Brewery      `url:"-"`
}

// LocationType represents the specific type of the Location.
type LocationType string

// Types of Locations.
const (
	LocationMicro      LocationType = "micro"
	LocationMacro                   = "macro"
	LocationNano                    = "nano"
	LocationBrewpub                 = "brewpub"
	LocationProduction              = "production"
	LocationOffice                  = "office"
	LocationTasting                 = "tasting"
	LocationRestaurant              = "restaurant"
	LocationCidery                  = "cidery"
	LocationMeadery                 = "meadery"
)

// LocationOrder specifies the ordering of a LocationList.
type LocationOrder string

// LocationList ordering options.
const (
	LocationOrderName           LocationOrder = "name"
	LocationOrderBreweryname                  = "breweryName"
	Locality                                  = "locality"
	LocationOrderRegion                       = "region"
	LocationOrderPostalCode                   = "postalCode"
	LocationOrderIsPrimary                    = "isPrimary"
	LocationOrderInPlanning                   = "inPlanning"
	LocationOrderIsClosed                     = "isClosed"
	LocationOrderLocationType                 = "locationType"
	LocationOrderCountryIsoCode               = "countryIsoCode"
	LocationOrderStatus                       = "status"
	LocationOrderCreateDate                   = "createDate"
	LocationOrderUpdateDate                   = "updateDate"
)

// LocationListRequest contains options for specifying Locations.
type LocationListRequest struct {
	Page           int           `url:"p,omitempty"`
	IDs            string        `url:"ids,omitempty"`
	Locality       string        `url:"locality,omitempty"`
	Region         string        `url:"region,omitempty"`
	PostalCode     string        `url:"postalCode,omitempty"`
	IsPrimary      YesNo         `url:"isPrimary,omitempty"`
	InPlanning     YesNo         `url:"inPlanning,omitempty"`
	IsClosed       YesNo         `url:"isClosed,omitempty"`
	LocationType   LocationType  `url:"locationType,omitempty"`
	CountryISOCode string        `url:"countryIsoCode,omitempty"`
	Since          int           `url:"since,omitempty"`
	Status         string        `url:"status,omitempty"`
	Order          LocationOrder `url:"order,omitempty"`
	Sort           ListSort      `url:"sort,omitempty"`
}

// LocationList represents a "page" containing a slice of Locations.
type LocationList struct {
	NumberOfPages int
	CurrentPage   int
	TotalResults  int
	Locations     []Location `json:"data"`
}

// List retrieves a list of Locations matching the given request.
// For non-premium users, one of Locality, PostalCode, Region must be set.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/location_index#1
func (ls *LocationService) List(q *LocationListRequest) (ll LocationList, err error) {
	// GET: /locations
	var req *http.Request
	req, err = ls.c.NewRequest("GET", "/locations", q)
	if err != nil {
		return
	}

	err = ls.c.Do(req, &ll)
	return
}

// Get retrieves the Location with the given ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/location_index#2
func (ls *LocationService) Get(locID string) (l Location, err error) {
	// GET: /location/:locationID
	var req *http.Request
	req, err = ls.c.NewRequest("GET", "/location/"+locID, nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    Location
		Message string
	}{}
	err = ls.c.Do(req, &resp)
	return resp.Data, err
}

// Update updates the Location having the given ID to match the given Location.
// The CountryISOCode of the given Location *must* be set.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/location_index#3
func (ls *LocationService) Update(locID string, l *Location) error {
	// PUT: /location/:locationID
	if l == nil {
		return fmt.Errorf("nil Location")
	}
	req, err := ls.c.NewRequest("PUT", "/location/"+locID, l)
	if err != nil {
		return err
	}

	// TODO: return any response?
	return ls.c.Do(req, nil)
}

// Delete removes the Location with the given ID from the BreweryDB.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/location_index#4
func (ls *LocationService) Delete(locID string) error {
	// DELETE: /location/:locationID
	req, err := ls.c.NewRequest("DELETE", "/location/"+locID, nil)
	if err != nil {
		return err
	}

	return ls.c.Do(req, nil)
}
