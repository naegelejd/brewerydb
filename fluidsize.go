package brewerydb

import (
	"fmt"
	"net/http"
)

// FluidsizeService provides access to the BreweryDB Fluidsize API.
// Use Client.Fluidsize.
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
func (fs *FluidsizeService) List() (fl []Fluidsize, err error) {
	// GET: /fluidsizes
	var req *http.Request
	req, err = fs.c.NewRequest("GET", "/fluidsizes", nil)
	if err != nil {
		return
	}

	fluidsizesResponse := struct {
		Status  string
		Data    []Fluidsize
		Message string
	}{}
	if err = fs.c.Do(req, &fluidsizesResponse); err != nil {
		return
	}

	return fluidsizesResponse.Data, nil
}

// Get returns the Fluidsize with the given Fluidsize ID.
func (fs *FluidsizeService) Get(id int) (f Fluidsize, err error) {
	// GET: /fluidsize/:fluidsizeId
	var req *http.Request
	req, err = fs.c.NewRequest("GET", fmt.Sprintf("/fluidsize/%d", id), nil)
	if err != nil {
		return
	}

	fluidsizeResponse := struct {
		Status  string
		Data    Fluidsize
		Message string
	}{}

	if err = fs.c.Do(req, &fluidsizeResponse); err != nil {
		return
	}

	return fluidsizeResponse.Data, nil
}
