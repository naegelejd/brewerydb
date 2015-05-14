package brewerydb

import (
	"fmt"
	"net/http"
)

// FeatureService provides access to the BreweryDB Feature API. Use Client.Feature.
type FeatureService struct {
	c *Client
}

// Feature represents a combined Featured Beer and Brewery.
// TODO: the Brewery in a Feature should ALSO contain its locations:
// see: http://www.brewerydb.com/developers/docs-endpoint/feature_index#1
type Feature struct {
	BeerID  string
	Beer    Beer
	Brewery Brewery
}

// Featured returns this week's Featured Beer and Brewery.
func (fs *FeatureService) Featured() (f Feature, err error) {
	// GET: /featured
	var req *http.Request
	req, err = fs.c.NewRequest("GET", "/featured", nil)
	if err != nil {
		return
	}

	featuredResponse := struct {
		Status  string
		Data    Feature
		Message string
	}{}

	if err = fs.c.Do(req, &featuredResponse); err != nil {
		return
	}
	return featuredResponse.Data, nil
}

// FeatureListRequest contains options for querying for a list of features.
type FeatureListRequest struct {
	Page         int    `url:"p"`
	Year         int    `url:"year,omitempty"`
	Week         int    `url:"week,omitempty"`
	IgnoreFuture string `url:"ignoreFuture,omitempty"` // Y or N
}

// FeatureList represents a single "page" containing a slice of Features.
type FeatureList struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Features      []Feature `json:"data"`
}

// List returns all Featured Beers and Breweries.
func (fs *FeatureService) List(q *FeatureListRequest) (fl FeatureList, err error) {
	// GET: /features
	var req *http.Request
	req, err = fs.c.NewRequest("GET", "/features", q)
	if err != nil {
		return
	}

	err = fs.c.Do(req, &fl)
	return
}

// FeatureByWeek returns the Featured Beer and Brewery for the given
// year and week number.
func (fs *FeatureService) FeatureByWeek(year, week int) (f Feature, err error) {
	// GET: /feature/:year-week
	var req *http.Request
	req, err = fs.c.NewRequest("GET", fmt.Sprintf("/feature/%4d-%d", year, week), nil)
	if err != nil {
		return
	}

	err = fs.c.Do(req, &f)
	return
}
