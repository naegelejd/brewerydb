package brewerydb

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestBreweryGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("brewery.get.json", t)
	defer data.Close()

	const id = "jmGoBA"
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, id)
		io.Copy(w, data)
	})

	b, err := client.Brewery.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if b.ID != id {
		t.Fatalf("Brewery ID = %v, want %v", b.ID, id)
	}
}

func TestBreweryList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("brewery.list.json", t)
	defer data.Close()

	const established = "1988"
	mux.HandleFunc("/breweries/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		if v := r.FormValue("established"); v != established {
			t.Fatalf("Request.FormValue established = %v, wanted %v", v, established)
		}
		// TODO: check more request query values
		io.Copy(w, data)
	})

	bl, err := client.Brewery.List(&BreweryListRequest{Established: established})
	if err != nil {
		t.Fatal(err)
	}
	if len(bl.Breweries) <= 0 {
		t.Fatal("Expected >0 breweries")
	}

	for _, b := range bl.Breweries {
		if l := 6; l != len(b.ID) {
			t.Fatalf("Brewery ID len = %d, wanted %d", len(b.ID), l)
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
