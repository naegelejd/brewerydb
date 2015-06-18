package brewerydb

import (
	"fmt"
	"net/http"
)

// FluidsizeService provides access to the BreweryDB Fluidsize API.
// Use Client.Fluidsize.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/fluidsize_index
type FluidsizeService struct {
	c *Client
}

// Volume represents a fluidsize volume.
type Volume string

// Pre-defined fluidsize Volumes.
const (
	VolumeBarrel Volume = "barrel"
	VolumePack          = "pack"
	VolumeOunce         = "oz"
	VolumeLiter         = "liter"
)

// Fluidsize represents a Fluidsize assigned to a UPC code.
type Fluidsize struct {
	ID            int
	Volume        string
	VolumeDisplay string
	Quantity      string
	CreateDate    string
}

// List returns a list of Fluidsizes.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/fluidsize_index#1
func (fs *FluidsizeService) List() (fl []Fluidsize, err error) {
	// GET: /fluidsizes
	var req *http.Request
	req, err = fs.c.NewRequest("GET", "/fluidsizes", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Fluidsize
		Message string
	}{}
	err = fs.c.Do(req, &resp)
	return resp.Data, err
}

// Get returns the Fluidsize with the given Fluidsize ID.
//
// See: http://www.brewerydb.com/developers/docs-endpoint/fluidsize_index#2
func (fs *FluidsizeService) Get(id int) (f Fluidsize, err error) {
	// GET: /fluidsize/:fluidsizeId
	var req *http.Request
	req, err = fs.c.NewRequest("GET", fmt.Sprintf("/fluidsize/%d", id), nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    Fluidsize
		Message string
	}{}
	err = fs.c.Do(req, &resp)
	return resp.Data, err
}
