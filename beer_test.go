package brewerydb

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestBeerList(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/beer.list.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	defer data.Close()

	mux.HandleFunc("/beers/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		io.Copy(w, data)
	})

	bl, err := client.Beer.List(&BeerListRequest{ABV: "8"})
	if err != nil {
		t.Fatal(err)
	}
	if len(bl.Beers) <= 0 {
		t.Fatal("Expected >0 beers")
	}
	for _, b := range bl.Beers {
		if l := 6; l != len(b.ID) {
			t.Fatal("Beer ID len = %d, wanted %d", len(b.ID), l)
		}
	}
}

func ExampleBeerService_List() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

	// Get first 40 beers with an ABV between 8.0 and 9.0, descending, alphabetical
	bl, err := c.Beer.List(&BeerListRequest{ABV: "8", Sort: SortDescending})
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range bl.Beers {
		fmt.Println(b.Name, b.ID)
	}
}

func ExampleBeerService_Breweries() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

	// Get breweries for a given beer
	breweries, err := c.Beer.Breweries("jmGoBA")
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range breweries {
		fmt.Println(b.Name)
	}
}
