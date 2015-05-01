// Copyright 2015 Joseph Naegele. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package brewerydb provides bindings to the BreweryDB API
// (http://www.brewerydb.com)
package brewerydb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

var apiURL = "http://api.brewerydb.com/v2"

// Client serves as the interface to the BreweryDB API.
type Client struct {
	http.Client
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

func (c *Client) getJSON(url string, data interface{}) error {
	resp, err := c.Get(url)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	// for debugging: dump json to stdout
	reader := io.TeeReader(resp.Body, os.Stdout)

	if err := json.NewDecoder(reader).Decode(data); err != nil {
		return err
	}

	return nil
}
