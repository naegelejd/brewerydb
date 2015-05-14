package brewerydb

// GeoPointUnit differentiates between miles and kilometers.
type GeoPointUnit string

// Units of measurement.
const (
	Miles      GeoPointUnit = "mi"
	Kilometers GeoPointUnit = "km"
)

// GeoPointRequest contains options for specifying a geographic coordinate.
type GeoPointRequest struct {
	Latitude           float64      `url:"lat"` // Required
	Longitude          float64      `url:"lng"` // Required
	Radius             float64      `url:"radius,omitempty"`
	Unit               GeoPointUnit `url:"unit,omitempty"`               // Default: mi
	WithSocialAccounts string       `url:"withSocialAccounts,omitempty"` // Y/N
	WithGuilds         string       `url:"withGuilds,omitempty"`         // Y/N
	WithAlternateNames string       `url:"withAlternateNames,omitempty"` // Y/N
}

// GeoPoint searches for Locations near the geographic coordinate specified in the GeoPointRequest.
// TODO: pagination??
func (ss *SearchService) GeoPoint(q *GeoPointRequest) ([]Location, error) {
	req, err := ss.c.NewRequest("GET", "/search/geo/point", q)
	if err != nil {
		return nil, err
	}

	geoPointResult := struct {
		NumberOfPages int
		CurrentPage   int
		TotalResults  int
		Data          []Location
	}{}

	err = ss.c.Do(req, &geoPointResult)
	return geoPointResult.Data, err
}
