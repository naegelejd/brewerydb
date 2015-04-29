package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ConvertIDService provides access to the BreweryDB ID Conversion API.
// Use Client.ConvertID.
type ConvertIDService struct {
	c *Client
}

// ConvertType is the type of ID to convert.
type ConvertType string

const (
	ConvertBrewery ConvertType = "brewery" // ConvertBrewery converts Brewery IDs.
	ConvertBeer                = "beer"    // ConvertBeer converts Beer IDs.
)

// ConvertIDs converts a series of "old" Beer or Brewery IDs to the "new" format
// (BreweryDB v1 to v2)
func (c *ConvertIDService) ConvertIDs(t ConvertType, oldIDs ...int) (map[int]string, error) {
	// POST: /convertid
	v := url.Values{}
	v.Set("type", string(t))

	var ids []string
	for _, id := range oldIDs {
		ids = append(ids, strconv.Itoa(id))
	}
	v.Set("ids", strings.Join(ids, ","))

	u := c.c.url("/convertid", nil)
	resp, err := c.c.PostForm(u, v)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to convert ids")
	}
	defer resp.Body.Close()

	convertidResponse := struct {
		Status string
		Data   []struct {
			OldID int
			NewID string
		}
		Message string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&convertidResponse); err != nil {
		return nil, err
	}

	m := make(map[int]string)
	for _, d := range convertidResponse.Data {
		m[d.OldID] = d.NewID
	}

	return m, nil
}