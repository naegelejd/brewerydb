package brewerydb

import (
	"fmt"
	"net/http"
)

// AdjunctService provides access to the BreweryDB Adjunct API. Use Client.Adjunct.
type AdjunctService struct {
	c *Client
}

// Adjunct represents an additional ingredient used in making a Beer.
type Adjunct struct {
	ID              int
	Name            string
	Category        string // This will always be set to "misc"
	CategoryDisplay string // This will always be set to "Miscellaneous"
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

	var req *http.Request
	req, err = as.c.NewRequest("GET", "/adjuncts", &Page{page})
	if err != nil {
		return
	}

	err = as.c.Do(req, &al)
	return
}

// Get obtains the Adjunct with the given Adjunct ID.
func (as *AdjunctService) Get(id int) (a Adjunct, err error) {
	// GET: /adjunct/:adjunctID
	var req *http.Request
	req, err = as.c.NewRequest("GET", fmt.Sprintf("/adjunct/%d", id), nil)
	if err != nil {
		return
	}

	adjunctReponse := struct {
		Status  string
		Data    Adjunct
		Message string
	}{}
	if err = as.c.Do(req, &adjunctReponse); err != nil {
		return
	}

	return adjunctReponse.Data, nil
}
