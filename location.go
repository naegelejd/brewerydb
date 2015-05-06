package brewerydb

import "net/http"

// LocationService provides access to the BreweryDB Location API.
// Use Client.Location.
type LocationService struct {
	c *Client
}

// Country represents a country found on Earth.
type Country struct {
	IsoCode     string
	Name        string
	DisplayName string
	IsoThree    string
	NumberCode  int
	URLTitle    string
	CreateDate  string
	// UpdateDate  string
}

// Location represents a the location of a Brewery or similar institution.
type Location struct {
	ID                       string
	Name                     string `json:"name,omitempty"`
	StreetAddress            string `json:"streetAddress,omitempty"`
	ExtendedAddress          string `json:"extendedAddress,omitempty"`
	Locality                 string `json:"locality,omitempty"`
	Region                   string `json:"region,omitempty"`
	PostalCode               string `json:"postalCode,omitempty"`
	Phone                    string `json:"phone,omitempty"`
	Website                  string `json:"website,omitempty"`
	HoursOfOperation         string
	HoursOfOperationExplicit []string     `json:"hoursOfOperationExplicit,omitempty"`
	HoursOfOperationNotes    string       `json:"hoursOfOperationNotes,omitempty"`
	TourInfo                 string       `json:"tourInfo,omitempty"`
	TimezoneID               string       `json:"timezoneId,omitempty"`
	Latitude                 float64      `json:"latitude,omitempty"`
	Longitude                float64      `json:"longitude,omitempty"`
	IsPrimary                string       `json:"isPrimary,omitempty"`    // Y/N
	InPlanning               string       `json:"inPlanning,omitempty"`   // Y/N
	IsClosed                 string       `json:"isClosed,omitempty"`     // Y/N
	OpenToPublic             string       `json:"openToPublic,omitempty"` // Y/N
	LocationType             LocationType `json:"locationType,omitempty"`
	LocationTypeDisplay      string
	CountryISOCode           string `json:"countryIsoCode"` // Required for UpdateLocation
	Country                  Country
	CreateDate               string
	UpdateDate               string
	YearOpened               string
	YearClosed               string
	BreweryID                string
	Brewery                  Brewery
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

// LocationRequest contains options for specifying Locations.
type LocationRequest struct {
	Page           int           `json:"p,omitempty"`
	IDs            string        `json:"ids,omitempty"`
	Locality       string        `json:"locality,omitempty"`
	Region         string        `json:"region,omitempty"`
	PostalCode     string        `json:"postalCode,omitempty"`
	IsPrimary      string        `json:"isPrimary,omitempty"`
	InPlanning     string        `json:"inPlanning,omitempty"`
	IsClosed       string        `json:"isClosed,omitempty"`
	LocationType   LocationType  `json:"locationType,omitempty"`
	CountryISOCode string        `json:"countryIsoCode,omitempty"`
	Since          int           `json:"since,omitempty"`
	Status         string        `json:"status,omitempty"`
	Order          LocationOrder `json:"order,omitempty"`
	Sort           ListSort      `json:"sort,omitempty"`
}

// LocationList represents a "page" containing a slice of Locations.
type LocationList struct {
	NumberOfPages int
	CurrentPage   int
	TotalResults  int
	Locations     []Location `json:"data"`
}

// List retrieves a list of Locations matching the given request.
func (ls *LocationService) List(q *LocationRequest) (ll LocationList, err error) {
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
func (ls *LocationService) Get(locID string) (l Location, err error) {
	// GET: /location/:locationID
	var req *http.Request
	req, err = ls.c.NewRequest("GET", "/location/"+locID, nil)
	if err != nil {
		return
	}

	locationResponse := struct {
		Status  string
		Data    Location
		Message string
	}{}

	err = ls.c.Do(req, &locationResponse)
	return locationResponse.Data, err
}

// UpdateLocation updates the Location having the given ID to match the given Location.
func (ls *LocationService) UpdateLocation(locID string, l *Location) error {
	// PUT: /location/:locationID
	req, err := ls.c.NewRequest("PUT", "/location/"+locID, l)
	if err != nil {
		return err
	}

	// TODO: return any response?
	return ls.c.Do(req, nil)
}

// DeleteLocation removes the Location with the given ID from the BreweryDB.
func (ls *LocationService) DeleteLocation(locID string) error {
	// DELETE: /location/:locationID
	req, err := ls.c.NewRequest("DELETE", "/location/"+locID, nil)
	if err != nil {
		return err
	}

	return ls.c.Do(req, nil)
}
