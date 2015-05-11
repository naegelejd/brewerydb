package brewerydb

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestBreweryList(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/brewery.list.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	defer data.Close()

	mux.HandleFunc("/breweries/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		io.Copy(w, data)
	})

	bl, err := client.Brewery.List(&BreweryListRequest{Established: "1988"})
	if err != nil {
		t.Fatal(err)
	}
	if len(bl.Breweries) <= 0 {
		t.Fatal("Expected >0 breweries")
	}

	for _, b := range bl.Breweries {
		if l := 6; l != len(b.ID) {
			t.Fatalf("Beer ID len = %d, wanted %d", len(b.ID), l)
		}
	}
}

func ExampleBreweryService_List() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

	// Get all breweries established in 1983
	bl, err := c.Brewery.List(&BreweryListRequest{Established: "1983"})
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range bl.Breweries {
		fmt.Println(b.Name, b.ID)
	}

	// Get all information about brewery with given ID (Flying Dog)
	b, err := c.Brewery.Get("jmGoBA")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(b.Name)
	fmt.Println(b.Description)
	fmt.Println(b.Website)
}
