package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

// Styles returns all Styles on the given page.
func (ss *StyleService) Styles(page int) (sl StyleList, err error) {
	// GET: /styles
	v := url.Values{}
	v.Set("p", strconv.Itoa(page))
	u := ss.c.url("/styles", &v)
	var resp *http.Response
	resp, err = ss.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get styles")
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&sl); err != nil {
		return
	}

	return
}

// Style obtains the Style with the given Style ID.
func (ss *StyleService) Style(id int) (s Style, err error) {
	// GET: /style/:styleID
	u := ss.c.url(fmt.Sprintf("/style/%d", id), nil)
	var resp *http.Response
	resp, err = ss.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get style")
		return
	}
	defer resp.Body.Close()

	styleResponse := struct {
		Status  string
		Data    Style
		Message string
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&styleResponse); err != nil {
		return
	}
	s = styleResponse.Data

	return
}
