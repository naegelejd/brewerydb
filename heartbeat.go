package brewerydb

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
	req, err := hs.c.NewRequest("GET", "/heartbeat", nil)
	if err != nil {
		return err
	}

	return hs.c.Do(req, nil)
}
