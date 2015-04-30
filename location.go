package brewerydb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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
	UrlTitle    string
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

// Locations retrieves a list of Locations matching the given request.
func (ls *LocationService) Locations(req *LocationRequest) (ll LocationList, err error) {
	// GET: /locations
	v := encode(req)
	u := ls.c.url("/locations", &v)

	var resp *http.Response
	resp, err = ls.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get locations")
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&ll); err != nil {
		return
	}

	return
}

// Location retrieves the Location with the given ID.
func (ls *LocationService) Location(locID string) (l Location, err error) {
	// GET: /location/:locationID
	u := ls.c.url("/location/"+locID, nil)

	var resp *http.Response
	resp, err = ls.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get location")
		return
	}
	defer resp.Body.Close()

	locationResponse := struct {
		Status  string
		Data    Location
		Message string
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&locationResponse); err != nil {
		return
	}
	l = locationResponse.Data

	return
}

// UpdateLocation updates the Location having the given ID to match the given Location.
func (ls *LocationService) UpdateLocation(locID string, l *Location) error {
	// PUT: /location/:locationID
	u := ls.c.url("/location/"+locID, nil)
	v := encode(l)
	put, err := http.NewRequest("PUT", u, bytes.NewBufferString(v.Encode()))
	if err != nil {
		return err
	}

	resp, err := ls.c.Do(put)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to update location")
	}
	defer resp.Body.Close()

	// TODO: return any response?
	return nil
}

// DeleteLocation removes the Location with the given ID from the BreweryDB.
func (ls *LocationService) DeleteLocation(locID string) error {
	// DELETE: /location/:locationID
	u := ls.c.url("/location/"+locID, nil)

	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	resp, err := ls.c.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to delete location")
	}
	defer resp.Body.Close()

	return nil
}
