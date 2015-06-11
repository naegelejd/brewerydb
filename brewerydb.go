// Copyright 2015 Joseph Naegele. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package brewerydb provides bindings to the BreweryDB API
// (http://www.brewerydb.com)
package brewerydb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// apiURL is not const so it can be stubbed in unit tests.
var apiURL = "http://api.brewerydb.com/v2"

// Page is a convenience type for encoding only a page number
// when paginating lists.
type Page struct {
	P int `url:"p"`
}

// Images is a collection of up to three differently-sized image URLs.
type Images struct {
	Icon   string `url:"-"`
	Medium string `url:"-"`
	Large  string `url:"-"`
}

// YesNo is just a bool that is url-encoded into either "Y" or "S".
type YesNo bool

// EncodeValues adds the value "Y" or "N" to the given url.Values
// for the given key if the YesNo value is true or false, respectively.
func (yn YesNo) EncodeValues(key string, v *url.Values) error {
	if yn {
		v.Set(key, "Y")
	} else {
		v.Set(key, "N")
	}
	return nil
}

// UnmarshalJSON decodes the JSON value "Y" or "N" into a boolean
// true or false, respectively.
func (yn *YesNo) UnmarshalJSON(data []byte) error {
	// expect a single-rune string containing either 'Y' or 'N'
	yes, no := []byte{'"', 'Y', '"'}, []byte{'"', 'N', '"'}
	if bytes.Equal(data, yes) {
		*yn = true
	} else if bytes.Equal(data, no) {
		*yn = false
	} else {
		return fmt.Errorf("invalid JSON value for YesNo (%v)", string(data))
	}
	return nil
}

// Client serves as the interface to the BreweryDB API.
type Client struct {
	client      http.Client
	apiKey      string
	NumRequests int
	JSONWriter  io.Writer
	Adjunct     *AdjunctService
	Beer        *BeerService
	Brewery     *BreweryService
	Category    *CategoryService
	Change      *ChangeService
	ConvertID   *ConvertIDService
	Event       *EventService
	Feature     *FeatureService
	Fermentable *FermentableService
	Fluidsize   *FluidsizeService
	Glass       *GlassService
	Guild       *GuildService
	Heartbeat   *HeartbeatService
	Hop         *HopService
	Ingredient  *IngredientService
	Location    *LocationService
	Menu        *MenuService
	Search      *SearchService
	SocialSite  *SocialSiteService
	Style       *StyleService
	Yeast       *YeastService
}

// NewClient creates a new BreweryDB Client using the given API key.
func NewClient(apiKey string) *Client {
	c := &Client{}
	c.apiKey = apiKey
	c.Adjunct = &AdjunctService{c}
	c.Beer = &BeerService{c}
	c.Brewery = &BreweryService{c}
	c.Category = &CategoryService{c}
	c.Change = &ChangeService{c}
	c.ConvertID = &ConvertIDService{c}
	c.Event = &EventService{c}
	c.Feature = &FeatureService{c}
	c.Fermentable = &FermentableService{c}
	c.Fluidsize = &FluidsizeService{c}
	c.Glass = &GlassService{c}
	c.Guild = &GuildService{c}
	c.Heartbeat = &HeartbeatService{c}
	c.Hop = &HopService{c}
	c.Ingredient = &IngredientService{c}
	c.Location = &LocationService{c}
	c.Menu = &MenuService{c}
	c.Search = &SearchService{c}
	c.SocialSite = &SocialSiteService{c}
	c.Style = &StyleService{c}
	c.Yeast = &YeastService{c}
	return c
}

// NewRequest creates a new http.Request with the given method,
// BreweryDB endpoint, and optionally a struct to be URL-encoded
// in the request.
func (c *Client) NewRequest(method string, endpoint string, data interface{}) (req *http.Request, err error) {
	var u *url.URL
	u, err = url.Parse(apiURL)
	if err != nil {
		return
	}
	u.Path += endpoint

	var dataVals url.Values
	if data != nil {
		dataVals, err = query.Values(data)
		if err != nil {
			return
		}
	} else {
		dataVals = url.Values{}
	}

	switch method {
	case "GET":
		fallthrough
	case "DELETE":
		dataVals.Set("key", c.apiKey)
		u.RawQuery = dataVals.Encode()
		req, err = http.NewRequest(method, u.String(), nil)
	case "POST":
		fallthrough
	case "PUT":
		q := url.Values{}
		q.Set("key", c.apiKey)
		u.RawQuery = q.Encode()

		payload := dataVals.Encode()
		body := bytes.NewBufferString(payload)
		req, err = http.NewRequest(method, u.String(), body)
		if err != nil {
			return
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(payload)))
	default:
		err = fmt.Errorf("Unknown HTTP method: %s", method)
	}

	return
}

// Do performs the given http.Request and optionally
// decodes the JSON response into the given data struct.
func (c *Client) Do(req *http.Request, data interface{}) error {
	// TODO: [DEBUGGING] fmt.Println(req.Method, req.URL)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// TODO: return a more useful error message
		return fmt.Errorf("HTTP Error %d", resp.StatusCode)
	}

	c.NumRequests++

	if data != nil {
		var body io.Reader
		// if the client has a JSONWriter, also dump JSON responses
		if c.JSONWriter != nil {
			body = io.TeeReader(resp.Body, c.JSONWriter)
		} else {
			body = resp.Body
		}

		err = json.NewDecoder(body).Decode(data)
	}

	return err
}
