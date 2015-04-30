package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// AdjunctService provides access to the BreweryDB Adjunct API. Use Client.Adjunct.
type AdjunctService struct {
	c *Client
}

// Adjunct represents an additional ingredient used in making a Beer.
type Adjunct struct {
	ID              int
	Name            string
	Category        string
	CategoryDisplay string
	CreateDate      string
}

// AdjunctList represents a single "page" containing a slice of Adjuncts.
type AdjunctList struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Adjuncts      []Adjunct `json:"data"`
}

// List returns all Adjuncts on the given page.
func (as *AdjunctService) List(page int) (al AdjunctList, err error) {
	// GET: /adjuncts
	v := url.Values{}
	v.Set("p", strconv.Itoa(page))
	u := as.c.url("/adjuncts", &v)
	var resp *http.Response
	resp, err = as.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get adjuncts")
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&al); err != nil {
		return
	}

	return
}

// Get obtains the Adjunct with the given Adjunct ID.
func (as *AdjunctService) Get(id int) (a Adjunct, err error) {
	// GET: /adjunct/:adjunctID
	u := as.c.url(fmt.Sprintf("/adjunct/%d", id), nil)
	var resp *http.Response
	resp, err = as.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get adjunct")
		return
	}
	defer resp.Body.Close()

	adjunctReponse := struct {
		Status  string
		Data    Adjunct
		Message string
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&adjunctReponse); err != nil {
		return
	}
	a = adjunctReponse.Data

	return
}
