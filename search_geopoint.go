package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GeoPointResult represents a single brewery near a geographic coordinate.
type GeoPointResult struct {
	ID                  string
	Name                string
	BreweryID           string
	LocationTypeDisplay string
	YearOpened          string
	YearClosed          string
	Distance            float64
	Latitude            float64
	Longitude           float64
	StreetAddress       string
	PostalCode          string
	Locality            string
	Region              string
	CountryISOCode      string
	Phone               string
	Status              string
	StatusDisplay       string
	CreateDate          string
	UpdateDate          string
	InPlanning          string // Y/N
	OpenToPublic        string // Y/N
	IsClosed            string // Y/N
	IsPrimary           string // Y/N
	Brewery             Brewery
}

// GeoPointUnit differentiates between miles and kilometers.
type GeoPointUnit string

const (
	Miles      GeoPointUnit = "mi"
	Kilometers GeoPointUnit = "km"
)

// GeoPointRequest contains options for specifying a geographic coordinate.
type GeoPointRequest struct {
	Latitude           float64      `json:"lat"` // Required
	Longitude          float64      `json:"lng"` // Required
	Radius             float64      `json:"radius,omitempty"`
	Unit               GeoPointUnit `json:"unit,omitempty"`               // Default: mi
	WithSocialAccounts string       `json:"withSocialAccounts,omitempty"` // Y/N
	WithGuilds         string       `json:"withGuilds,omitempty"`         // Y/N
	WithAlternateNames string       `json:"withAlternateNames,omitempty"` // Y/N
}

// GeoPoint searches for Breweries near the geographic coordinate specified in the GeoPointRequest.
func (ss *SearchService) GeoPoint(req *GeoPointRequest) ([]GeoPointResult, error) {
	v := encode(req)
	u := ss.c.url("/search/geo/point", &v)
	resp, err := ss.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to search geo point")
	}
	defer resp.Body.Close()

	geoPointResult := struct {
		NumberOfPages int
		CurrentPage   int
		TotalResults  int
		Data          []GeoPointResult
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&geoPointResult); err != nil {
		return nil, err
	}

	return geoPointResult.Data, nil
}
