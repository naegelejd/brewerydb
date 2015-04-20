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

type Client struct {
	http.Client
	apiKey      string
	numRequests int
}

func NewClient(apiKey string) *Client {
	return &Client{http.Client{}, apiKey, 0}
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
