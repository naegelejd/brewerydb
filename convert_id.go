package brewerydb

import (
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

// Types of things that have IDs that can be converted.
const (
	ConvertBrewery ConvertType = "brewery" // ConvertBrewery converts Brewery IDs.
	ConvertBeer                = "beer"    // ConvertBeer converts Beer IDs.
)

// ConvertIDs converts a series of "old" Beer or Brewery IDs to the "new" format
// (BreweryDB v1 to v2)
func (cs *ConvertIDService) ConvertIDs(t ConvertType, oldIDs ...int) (map[int]string, error) {
	// POST: /convertid

	var ids []string
	for _, id := range oldIDs {
		ids = append(ids, strconv.Itoa(id))
	}

	q := struct {
		Type string `url:"type"`
		IDs  string `url:"ids"`
	}{string(t), strings.Join(ids, ",")}

	req, err := cs.c.NewRequest("POST", "/convertid", &q)
	if err != nil {
		return nil, err
	}

	resp := struct {
		Status string
		Data   []struct {
			OldID int
			NewID string
		}
		Message string
	}{}
	if err := cs.c.Do(req, &resp); err != nil {
		return nil, err
	}

	m := make(map[int]string)
	for _, d := range resp.Data {
		m[d.OldID] = d.NewID
	}

	return m, nil
}
