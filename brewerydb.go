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
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const apiURL = "http://api.brewerydb.com/v2"

// Page is a convenience type for encoding only a page number
// when paginating lists.
type Page struct {
	P int `json:"p"`
}

// Client serves as the interface to the BreweryDB API.
type Client struct {
	client      http.Client
	apiKey      string
	numRequests int
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
func (c *Client) NewRequest(method string, endpoint string, data interface{}) (*http.Request, error) {
	url := apiURL + endpoint + "/?key=" + c.apiKey
	var body io.Reader
	if data != nil {
		vals := encode(data)
		payload := vals.Encode()
		if method == "GET" {
			url += "&" + payload
		} else {
			body = bytes.NewBufferString(payload)
		}
	}

	// debugging:
	log.Println(url)

	return http.NewRequest(method, url, body)
}

// Do performs the given http.Request and optionally
// decodes the JSON response into the given data struct.
func (c *Client) Do(req *http.Request, data interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// TODO: return a more useful error message
		return fmt.Errorf("HTTP Error %d", resp.StatusCode)
	}

	if data != nil {
		// debugging:
		body := io.TeeReader(resp.Body, os.Stdout)
		// body := resp.Body
		if w, ok := data.(io.Writer); ok {
			_, err = io.Copy(w, body)
		} else {
			err = json.NewDecoder(body).Decode(data)
		}
	}

	return err
}

// TODO: DELETE THIS
func (c *Client) url(endpoint string, vals *url.Values) string {
	if vals == nil {
		vals = &url.Values{}
	}
	vals.Set("key", c.apiKey)
	query := vals.Encode()
	u := apiURL + endpoint + "/?" + query
	fmt.Println(u)
	return u
}

// TODO: DELETE THIS
func (c *Client) Get(url string) (*http.Response, error) {
	return c.client.Get(url)
}
