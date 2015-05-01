package brewerydb

import (
	"encoding/json"
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
	u := fs.c.url("/featured", nil)

	var resp *http.Response
	resp, err = fs.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode, resp.Body)
		err = fmt.Errorf("unable to get featured")
		return
	}
	defer resp.Body.Close()

	featuredResponse := struct {
		Status  string
		Data    Feature
		Message string
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&featuredResponse); err != nil {
		return
	}
	f = featuredResponse.Data
	return
}

// FeatureRequest contains options for querying for a list of features.
type FeatureRequest struct {
	Page         int    `json:"p"`
	Year         int    `json:"year"`
	Week         int    `json:"week"`
	IgnoreFuture string `json:"ignoreFuture"` // Y or N
}

// FeatureList represents a single "page" containing a slice of Features.
type FeatureList struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Features      []Feature `json:"data"`
}

// List returns all Featured Beers and Breweries.
func (fs *FeatureService) List(req *FeatureRequest) (fl FeatureList, err error) {
	// GET: /features
	v := encode(req)
	u := fs.c.url("/features", &v)

	var resp *http.Response
	resp, err = fs.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get features")
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&fl); err != nil {
		return
	}
	return
}

// FeatureByWeek returns the Featured Beer and Brewery for the given
// year and week number.
func (fs *FeatureService) FeatureByWeek(year, week int) (f Feature, err error) {
	// GET: /feature/:year-week
	u := fs.c.url(fmt.Sprintf("/feature/%4d-%d", year, week), nil)

	var resp *http.Response
	resp, err = fs.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get feature")
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return
	}
	return
}
