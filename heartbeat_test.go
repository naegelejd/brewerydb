package brewerydb

import (
	"fmt"
	"net/http"
	"testing"
)

var fakeDataHeartbeat = `{
  "status" : "success",
  "data" : {
    "format" : "json",
    "requestMethod" : "GET",
    "key" : "",
    "timestamp" : 1337870979
  },
  "message" : "Request Successful"
}`

func TestHearbeat(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fakeDataHeartbeat)
	})

	if err := client.Heartbeat.Heartbeat(); err != nil {
		t.Fatal(err)
	}

	testBadURL(t, func() error {
		return client.Heartbeat.Heartbeat()
	})
}
