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

	mux.HandleFunc("/locations/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		io.Copy(w, data)
	})

	ll, err := client.Location.List(&LocationRequest{Page: 1, Region: "Maryland"})
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
