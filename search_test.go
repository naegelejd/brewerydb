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

	mux.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		io.Copy(w, data)
	})

	bl, err := client.Search.Beer("flying", &SearchRequest{Page: 1})
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

	mux.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		io.Copy(w, data)
	})

	bl, err := client.Search.Brewery("dog", &SearchRequest{Page: 1})
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
