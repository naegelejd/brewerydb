package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// FluidsizeService provides access to the BreweryDB Fluidsize API.
// Use Client.Fluidsize.
type FluidsizeService struct {
	c *Client
}

// Fluidsize represents a Fluidsize assigned to a UPC code.
type Fluidsize struct {
	ID            int
	Volume        string
	VolumeDisplay string
	Quantity      string
	CreateDate    string
}

// Fluidsizes returns a list of Fluidsizes.
func (fs *FluidsizeService) Fluidsizes() (fl []Fluidsize, err error) {
	// GET: /fluidsizes
	u := fs.c.url("/fluidsizes", nil)

	var resp *http.Response
	resp, err = fs.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get fluidsizes")
		return
	}
	defer resp.Body.Close()

	fluidsizesResponse := struct {
		Status  string
		Data    []Fluidsize
		Message string
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&fluidsizesResponse); err != nil {
		return
	}
	fl = fluidsizesResponse.Data
	return
}

// Fluidsize returns the Fluidsize with the given Fluidsize ID.
func (fs *FluidsizeService) Fluidsize(id int) (f Fluidsize, err error) {
	// GET: /fluidsize/:fluidsizeId
	u := fs.c.url(fmt.Sprintf("/fluidsize/%d", id), nil)

	var resp *http.Response
	resp, err = fs.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get fluidsize")
		return
	}
	defer resp.Body.Close()

	fluidsizeResponse := struct {
		Status  string
		Data    Fluidsize
		Message string
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&fluidsizeResponse); err != nil {
		return
	}
	f = fluidsizeResponse.Data
	return
}