package brewerydb

import (
	"fmt"
	"net/http"
)

// StyleService provides access to the BreweryDB Style API. Use Client.Style.
type StyleService struct {
	c *Client
}

// Style represents a style of Beer.
type Style struct {
	ID          int
	Name        string
	ShortName   string
	Description string
	CategoryID  int
	Category    Category
	IbuMin      string
	IbuMax      string
	SrmMin      string
	SrmMax      string
	OgMin       string
	OgMax       string
	FgMin       string
	FgMax       string
	AbvMin      string
	AbvMax      string
	CreateDate  string
	UpdateDate  string
}

// StyleList represents a single "page" containing a slice of Styles.
type StyleList struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Styles        []Style `json:"data"`
}

// List returns all Styles on the given page.
func (ss *StyleService) List(page int) (sl StyleList, err error) {
	// GET: /styles
	var req *http.Request
	req, err = ss.c.NewRequest("GET", "/styles", &Page{page})
	if err != nil {
		return
	}

	err = ss.c.Do(req, &sl)
	return
}

// Get obtains the Style with the given Style ID.
func (ss *StyleService) Get(id int) (s Style, err error) {
	// GET: /style/:styleID
	var req *http.Request
	req, err = ss.c.NewRequest("GET", fmt.Sprintf("/style/%d", id), nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    Style
		Message string
	}{}
	err = ss.c.Do(req, &resp)
	return resp.Data, err
}
