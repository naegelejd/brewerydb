package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// YeastService provides access to the BreweryDB Yeast API. Use Client.Yeast.
type YeastService struct {
	c *Client
}

// Yeast represents a type of yeast used in making a Beer.
type Yeast struct {
	ID              int
	Name            string
	Category        string
	CategoryDisplay string
	CreateDate      string
}

// YeastList represents a single "page" containing a slice of Yeasts.
type YeastList struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Yeasts        []Yeast `json:"data"`
}

// List returns all Yeasts on the given page.
func (ys *YeastService) List(page int) (yl YeastList, err error) {
	// GET: /yeasts
	v := url.Values{}
	v.Set("p", strconv.Itoa(page))
	u := ys.c.url("/yeasts", &v)
	var resp *http.Response
	resp, err = ys.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get yeasts")
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&yl); err != nil {
		return
	}

	return
}

// Get obtains the Yeast with the given Yeast ID.
func (ys *YeastService) Get(id int) (y Yeast, err error) {
	// GET: /yeast/:yeastID
	u := ys.c.url(fmt.Sprintf("/yeast/%d", id), nil)
	var resp *http.Response
	resp, err = ys.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get yeast")
		return
	}
	defer resp.Body.Close()

	yeastResponse := struct {
		Status  string
		Data    Yeast
		Message string
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&yeastResponse); err != nil {
		return
	}
	y = yeastResponse.Data

	return
}
