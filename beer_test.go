package brewerydb

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestBeerGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("beer.get.json", t)
	defer data.Close()

	const id = "o9TSOv"
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, id)
		io.Copy(w, data)
	})

	b, err := client.Beer.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if b.ID != id {
		t.Fatalf("Beer ID = %v, want %v", b.ID, id)
	}
}

func TestBeerList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("beer.list.json", t)
	defer data.Close()

	const abv = "8"
	mux.HandleFunc("/beers/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		abv := r.FormValue("abv")
		if v := r.FormValue("abv"); v != abv {
			t.Fatalf("Request.FormValue abv = %v, wanted %v", v, abv)
		}
		// TODO: check more request query values
		io.Copy(w, data)
	})

	bl, err := client.Beer.List(&BeerListRequest{ABV: abv})
	if err != nil {
		t.Fatal(err)
	}
	if len(bl.Beers) <= 0 {
		t.Fatal("Expected >0 beers")
	}
	for _, b := range bl.Beers {
		if l := 6; l != len(b.ID) {
			t.Fatalf("Beer ID len = %d, wanted %d", len(b.ID), l)
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
