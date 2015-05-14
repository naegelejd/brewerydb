package brewerydb

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestLocationList(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/location.list.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	defer data.Close()

	const (
		page   = 1
		region = "Maryland"
	)
	mux.HandleFunc("/locations/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)
		if v := r.FormValue("region"); v != region {
			t.Fatalf("Request.FormValue region = %v, want %v", v, region)
		}
		// TODO: check more request query values
		io.Copy(w, data)
	})

	ll, err := client.Location.List(&LocationListRequest{Page: page, Region: region})
	if err != nil {
		t.Fatal(err)
	}
	if len(ll.Locations) <= 0 {
		t.Fatal("Expected >0 locations")
	}

	for _, l := range ll.Locations {
		if n := 6; n != len(l.ID) {
			t.Fatalf("Location ID len = %d, wanted %d", len(l.ID), n)
		}
		if l.Latitude == 0.0 {
			t.Fatal("Expected non-zero latitude")
		}
		if l.Longitude == 0.0 {
			t.Fatal("Expected non-zero longitude")
		}
	}
}
