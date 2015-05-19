package brewerydb

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestLocationGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("location.get.json", t)
	defer data.Close()

	const id = "z9H6HJ"
	mux.HandleFunc("/location/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, id)
		io.Copy(w, data)
	})

	l, err := client.Location.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if l.ID != id {
		t.Fatalf("Location ID = %v, want %v", l.ID, id)
	}
}

func TestLocationList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("location.list.json", t)
	defer data.Close()

	const (
		page   = 1
		region = "Maryland"
	)
	mux.HandleFunc("/locations", func(w http.ResponseWriter, r *http.Request) {
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

func TestLocationDelete(t *testing.T) {
	setup()
	defer teardown()

	const id = "z9H6HJ"
	mux.HandleFunc("/location/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "location" {
			t.Fatal("bad URL, expected \"/location/:locationId\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Location ID", http.StatusNotFound)
		}

	})

	if err := client.Location.Delete(id); err != nil {
		t.Fatal(err)
	}

	if err := client.Location.Delete("******"); err == nil {
		t.Fatal("expected HTTP 404")
	}
}
