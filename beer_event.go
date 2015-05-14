package brewerydb

import "net/http"

// ListEvents returns a slice of Events where the given Beer is/was present
// or has won awards.
func (bs *BeerService) ListEvents(beerID string, onlyWinners bool) (el []Event, err error) {
	// GET: /beer/:beerId/events
	var q struct {
		OnlyWinners string `url:"onlyWinners,omitempty"`
	}
	if onlyWinners {
		q.OnlyWinners = "Y"
	}

	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+beerID+"/events", &q)
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
