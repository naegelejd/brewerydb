package brewerydb

import "net/http"

// ListEvents returns a slice of Events where the given Brewery is/was present
// or has won awards.
func (bs *BreweryService) ListEvents(breweryID string, onlyWinners bool) (el []Event, err error) {
	// GET: /brewery/:breweryId/events
	var q struct {
		OnlyWinners string `url:"onlyWinners,omitempty"`
	}
	if onlyWinners {
		q.OnlyWinners = "Y"
	}

	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/brewery/"+breweryID+"/events", &q)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Event
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}
