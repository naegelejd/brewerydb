package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// UPC retrieves one or more Beers matching the given Universal Product Code.
// TODO: pagination??
func (ss *SearchService) UPC(code string) ([]Beer, error) {
	v := url.Values{}
	v.Set("code", code)
	u := ss.c.url("/search/upc", &v)
	resp, err := ss.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to search UPC")
	}
	defer resp.Body.Close()

	upcResponse := struct {
		NumberOfPages int
		CurrentPage   int
		TotalResults  int
		Data          []Beer
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&upcResponse); err != nil {
		return nil, err
	}
	return upcResponse.Data, nil
}
