package brewerydb

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestSearchBeer(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/search.beer.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	defer data.Close()

	const (
		query = "flying"
		page  = 1
	)
	mux.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)
		if q := r.FormValue("q"); q != query {
			t.Fatalf("Request.FormValue q = %v, want %v", q, query)
		}
		// TODO: check more request query values
		io.Copy(w, data)
	})

	bl, err := client.Search.Beer(query, &SearchRequest{Page: page})
	if err != nil {
		t.Fatal(err)
	}
	if len(bl.Beers) <= 0 {
		t.Fatal("Expected >0 beers")
	}
	for _, b := range bl.Beers {
		if l := 6; l != len(b.ID) {
			t.Fatalf("Beer ID len = %d, want %d", len(b.ID), l)
		}
	}
}

func TestSearchBrewery(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/search.brewery.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	defer data.Close()

	const (
		query = "dog"
		page  = 1
	)
	mux.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)
		if q := r.FormValue("q"); q != query {
			t.Fatalf("Request.FormValue q = %v, want %v", q, query)
		}
		// TODO: check more request query values
		io.Copy(w, data)
	})

	bl, err := client.Search.Brewery(query, &SearchRequest{Page: page})
	if err != nil {
		t.Fatal(err)
	}
	if len(bl.Breweries) <= 0 {
		t.Fatal("Expected >0 breweries")
	}
	for _, b := range bl.Breweries {
		if l := 6; l != len(b.ID) {
			t.Fatalf("Brewery ID len = %d, want %d", len(b.ID), l)
		}
	}
}
