package brewerydb

import (
	"fmt"
	"net/http"
)

// HeartbeatService provides access to the BreweryDB Heartbeat API.
// Use Client.Heartbeat.
type HeartbeatService struct {
	c *Client
}

// HeartbeatResponse represents a BreweryDB Heartbeat (essentially an "echo"
// of the Heartbeat request).
type HeartbeatResponse struct {
	Status string
	Data   struct {
		Format        string
		RequestMethod string
		Key           string
		Timestamp     int
	}
	Message string
}

// Heartbeat checks whether the BreweryDB API is currently active. It
// returns nil if the API is available and an error otherwise.
func (hs *HeartbeatService) Heartbeat() error {
	u := hs.c.url("/heartbeat", nil)
	resp, err := hs.c.Get(u)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to detect heartbeat")
	}
	defer resp.Body.Close()

	// TODO: do anything with response?

	return nil
}
