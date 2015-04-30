package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Style retrieves one or more Styles matching the given query string.
// TODO: pagination??
func (ss *SearchService) Style(query string, withDescriptions bool) ([]Style, error) {
	v := url.Values{}
	v.Set("q", query)
	if withDescriptions {
		v.Set("withDescriptions", "Y")
	}
	u := ss.c.url("/search/style", &v)
	resp, err := ss.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to search style")
	}
	defer resp.Body.Close()

	styleResponse := struct {
		NumberOfPages int
		CurrentPage   int
		TotalResults  int
		Data          []Style
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&styleResponse); err != nil {
		return nil, err
	}
	return styleResponse.Data, nil
}
