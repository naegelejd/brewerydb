package brewerydb

import (
	"net/http"
	"net/url"
)

type Client struct {
	http.Client
	apiKey      string
	numRequests int
}

func NewClient(apiKey string) *Client {
	return &Client{http.Client{}, apiKey, 0}
}

func (c *Client) URL(endpoint string, vals *url.Values) string {
	if vals == nil {
		vals = &url.Values{}
	}
	vals.Set("key", c.apiKey)
	query := vals.Encode()
	return "http://api.brewerydb.com/v2" + endpoint + "/?" + query
}
